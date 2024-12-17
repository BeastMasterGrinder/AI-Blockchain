package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Point represents a 2D point.
type Point struct {
	X, Y float64
}

// Cluster represents a cluster with a centroid and its points.
type Cluster struct {
	Centroid Point
	Points   []Point
}

// Distance calculates the Euclidean distance between two points.
func Distance(a, b Point) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

// AssignPoints assigns each point to the nearest cluster based on the centroid.
func AssignPoints(points []Point, clusters []Cluster) {
	for i := range clusters {
		clusters[i].Points = nil // Clear existing points.
	}
	for _, p := range points {
		closestIndex := 0
		minDistance := Distance(p, clusters[0].Centroid)
		for i := 1; i < len(clusters); i++ {
			dist := Distance(p, clusters[i].Centroid)
			if dist < minDistance {
				minDistance = dist
				closestIndex = i
			}
		}
		clusters[closestIndex].Points = append(clusters[closestIndex].Points, p)
	}
}

// UpdateCentroids updates the centroids of the clusters based on the mean of their points.
func UpdateCentroids(clusters []Cluster) {
	for i := range clusters {
		if len(clusters[i].Points) == 0 {
			continue
		}
		var sumX, sumY float64
		for _, p := range clusters[i].Points {
			sumX += p.X
			sumY += p.Y
		}
		clusters[i].Centroid = Point{X: sumX / float64(len(clusters[i].Points)), Y: sumY / float64(len(clusters[i].Points))}
	}
}

// KMeans performs the K-Means clustering algorithm.
func KMeans(points []Point, k int, maxIterations int) []Cluster {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Initialize clusters with random centroids
	clusters := make([]Cluster, k)
	for i := range clusters {
		clusters[i].Centroid = points[rand.Intn(len(points))]
	}

	// Run the algorithm
	for i := 0; i < maxIterations; i++ {
		AssignPoints(points, clusters)
		UpdateCentroids(clusters)
	}

	return clusters
}

// ReadPoints reads 2D points from a file.
func ReadPoints(filePath string) ([]Point, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var points []Point
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		coords := strings.Split(line, ",")
		if len(coords) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}
		x, err := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse X coordinate: %v", err)
		}
		y, err := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Y coordinate: %v", err)
		}
		points = append(points, Point{X: x, Y: y})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return points, nil
}

// WriteClusters writes clustering results to a file.
func WriteClusters(filePath string, clusters []Cluster) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i, cluster := range clusters {
		_, err := fmt.Fprintf(writer, "Cluster %d:\n", i+1)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(writer, "  Centroid: (%.2f, %.2f)\n", cluster.Centroid.X, cluster.Centroid.Y)
		if err != nil {
			return err
		}
		_, err = writer.WriteString("  Points:\n")
		if err != nil {
			return err
		}
		for _, p := range cluster.Points {
			_, err = fmt.Fprintf(writer, "    (%.2f, %.2f)\n", p.X, p.Y)
			if err != nil {
				return err
			}
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func main() {
	// Input and output file paths
	inputFile := "input.txt"
	outputFile := "output.txt"

	// Read points from input file
	points, err := ReadPoints(inputFile)
	if err != nil {
		fmt.Printf("Error reading points: %v\n", err)
		return
	}

	// Number of clusters
	k := 3

	// Maximum number of iterations
	maxIterations := 100

	// Run K-Means
	clusters := KMeans(points, k, maxIterations)

	// Write clustering results to output file
	err = WriteClusters(outputFile, clusters)
	if err != nil {
		fmt.Printf("Error writing clusters: %v\n", err)
		return
	}

	fmt.Printf("Clustering results written to %s\n", outputFile)
}
