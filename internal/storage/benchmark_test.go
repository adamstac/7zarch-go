package storage

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// BenchmarkRegistryOperations tests performance of core registry operations
func BenchmarkRegistryOperations(b *testing.B) {
	registry, tempDir := setupTestRegistry(&testing.T{})
	defer func() {
		registry.Close()
		os.RemoveAll(tempDir)
	}()
	
	// Pre-populate registry with test data
	numArchives := 1000
	archives := make([]*Archive, numArchives)
	
	for i := 0; i < numArchives; i++ {
		archives[i] = &Archive{
			UID:      generateUID(),
			Name:     fmt.Sprintf("benchmark-archive-%04d.7z", i),
			Path:     fmt.Sprintf("/benchmark/archive-%04d.7z", i),
			Size:     int64(1024 * (i + 1)),
			Created:  time.Now().Add(-time.Duration(i) * time.Hour),
			Checksum: fmt.Sprintf("benchmark-checksum-%064d", i),
			Profile:  []string{"balanced", "media", "documents"}[i%3],
			Managed:  i%2 == 0,
			Status:   "present",
		}
		
		err := registry.Add(archives[i])
		if err != nil {
			b.Fatalf("Failed to add benchmark archive %d: %v", i, err)
		}
	}
	
	b.Run("Add", func(b *testing.B) {
		tempRegistry, tempDir := setupTestRegistry(&testing.T{})
		defer func() {
			tempRegistry.Close()
			os.RemoveAll(tempDir)
		}()
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			archive := &Archive{
				UID:     generateUID(),
				Name:    fmt.Sprintf("add-bench-%d.7z", i),
				Path:    fmt.Sprintf("/bench/add-bench-%d.7z", i),
				Size:    int64(1024 * i),
				Created: time.Now(),
				Profile: "balanced",
				Managed: true,
				Status:  "present",
			}
			
			err := tempRegistry.Add(archive)
			if err != nil {
				b.Fatalf("Add failed: %v", err)
			}
		}
	})
	
	b.Run("Get", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			name := fmt.Sprintf("benchmark-archive-%04d.7z", i%numArchives)
			_, err := registry.Get(name)
			if err != nil {
				b.Fatalf("Get failed: %v", err)
			}
		}
	})
	
	b.Run("List", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := registry.List()
			if err != nil {
				b.Fatalf("List failed: %v", err)
			}
		}
	})
	
	b.Run("ListNotUploaded", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := registry.ListNotUploaded()
			if err != nil {
				b.Fatalf("ListNotUploaded failed: %v", err)
			}
		}
	})
	
	b.Run("ListOlderThan", func(b *testing.B) {
		cutoff := 24 * time.Hour // 1 day
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := registry.ListOlderThan(cutoff)
			if err != nil {
				b.Fatalf("ListOlderThan failed: %v", err)
			}
		}
	})
}

// BenchmarkResolverOperations tests performance of AC's resolver (when implemented)
func BenchmarkResolverOperations(b *testing.B) {
	registry, tempDir := setupTestRegistry(&testing.T{})
	defer func() {
		registry.Close()
		os.RemoveAll(tempDir)
	}()
	
	// Create test data with known UIDs and checksums
	numArchives := 1000
	testUIDs := make([]string, numArchives)
	testChecksums := make([]string, numArchives)
	testNames := make([]string, numArchives)
	
	for i := 0; i < numArchives; i++ {
		uid := generateUID()
		checksum := fmt.Sprintf("checksum-%064d", i)
		name := fmt.Sprintf("resolve-test-%04d.7z", i)
		
		testUIDs[i] = uid
		testChecksums[i] = checksum
		testNames[i] = name
		
		archive := &Archive{
			UID:      uid,
			Name:     name,
			Path:     fmt.Sprintf("/resolve/test-%04d.7z", i),
			Size:     int64(1024 * (i + 1)),
			Created:  time.Now(),
			Checksum: checksum,
			Profile:  "balanced",
			Managed:  true,
			Status:   "present",
		}
		
		err := registry.Add(archive)
		if err != nil {
			b.Fatalf("Failed to add resolve test archive %d: %v", i, err)
		}
	}
	
	resolver := NewResolver(registry)
	
	b.Run("ResolveExactUID", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			uid := testUIDs[i%numArchives]
			_, err := resolver.Resolve(uid)
			if err != nil {
				b.Fatalf("Resolve exact UID failed: %v", err)
			}
		}
	})
	
	b.Run("ResolveUIDPrefix", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			uid := testUIDs[i%numArchives]
			prefix := uid[:8] // Use 8-character prefix
			_, err := resolver.Resolve(prefix)
			if err != nil {
				// May be ambiguous - that's expected
				if _, isAmbiguous := err.(*AmbiguousIDError); !isAmbiguous {
					b.Fatalf("Resolve UID prefix failed: %v", err)
				}
			}
		}
	})
	
	b.Run("ResolveChecksumPrefix", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			checksum := testChecksums[i%numArchives]
			prefix := checksum[:16] // Use 16-character prefix
			_, err := resolver.Resolve(prefix)
			if err != nil {
				if _, isAmbiguous := err.(*AmbiguousIDError); !isAmbiguous {
					b.Fatalf("Resolve checksum prefix failed: %v", err)
				}
			}
		}
	})
	
	b.Run("ResolveName", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			name := testNames[i%numArchives]
			_, err := resolver.Resolve(name)
			if err != nil {
				b.Fatalf("Resolve name failed: %v", err)
			}
		}
	})
}

// BenchmarkManagerOperations tests performance of storage manager
func BenchmarkManagerOperations(b *testing.B) {
	manager, tempDir := setupTestManager(&testing.T{})
	defer os.RemoveAll(tempDir)
	
	b.Run("ManagerAdd", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			name := fmt.Sprintf("manager-bench-%d.7z", i)
			path := fmt.Sprintf("/manager/bench-%d.7z", i)
			err := manager.Add(name, path, int64(1024*i), "balanced", "", "", true)
			if err != nil {
				b.Fatalf("Manager Add failed: %v", err)
			}
		}
	})
	
	b.Run("ManagerList", func(b *testing.B) {
		// Pre-populate with some data
		for i := 0; i < 100; i++ {
			name := fmt.Sprintf("pre-populate-%d.7z", i)
			path := fmt.Sprintf("/pre/populate-%d.7z", i)
			manager.Add(name, path, int64(1024*i), "balanced", "", "", true)
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := manager.List()
			if err != nil {
				b.Fatalf("Manager List failed: %v", err)
			}
		}
	})
	
	b.Run("ManagerGet", func(b *testing.B) {
		// Pre-populate with test data
		testNames := make([]string, 100)
		for i := 0; i < 100; i++ {
			name := fmt.Sprintf("get-test-%d.7z", i)
			path := fmt.Sprintf("/get/test-%d.7z", i)
			testNames[i] = name
			manager.Add(name, path, int64(1024*i), "balanced", "", "", true)
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			name := testNames[i%len(testNames)]
			_, err := manager.Get(name)
			if err != nil {
				b.Fatalf("Manager Get failed: %v", err)
			}
		}
	})
	
	b.Run("ManagerMarkUploaded", func(b *testing.B) {
		// Pre-populate with test data
		testNames := make([]string, 100)
		for i := 0; i < 100; i++ {
			name := fmt.Sprintf("upload-test-%d.7z", i)
			path := fmt.Sprintf("/upload/test-%d.7z", i)
			testNames[i] = name
			manager.Add(name, path, int64(1024*i), "balanced", "", "", true)
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			name := testNames[i%len(testNames)]
			destination := fmt.Sprintf("s3://benchmark/%s", name)
			err := manager.MarkUploaded(name, destination)
			if err != nil {
				b.Fatalf("Manager MarkUploaded failed: %v", err)
			}
		}
	})
}

// BenchmarkScalability tests performance with large datasets
func BenchmarkScalability(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping scalability benchmarks in short mode")
	}
	
	datasets := []int{100, 1000, 10000}
	
	for _, size := range datasets {
		b.Run(fmt.Sprintf("Registry%d", size), func(b *testing.B) {
			registry, tempDir := setupTestRegistry(&testing.T{})
			defer func() {
				registry.Close()
				os.RemoveAll(tempDir)
			}()
			
			// Populate registry with test data
			for i := 0; i < size; i++ {
				archive := &Archive{
					UID:     generateUID(),
					Name:    fmt.Sprintf("scale-test-%06d.7z", i),
					Path:    fmt.Sprintf("/scale/test-%06d.7z", i),
					Size:    int64(1024 * (i + 1)),
					Created: time.Now(),
					Profile: "balanced",
					Managed: true,
					Status:  "present",
				}
				
				err := registry.Add(archive)
				if err != nil {
					b.Fatalf("Failed to populate registry: %v", err)
				}
			}
			
			// Benchmark list operation with large dataset
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := registry.List()
				if err != nil {
					b.Fatalf("List failed with %d archives: %v", size, err)
				}
			}
		})
	}
}

// BenchmarkConcurrentAccess tests performance under concurrent load
func BenchmarkConcurrentAccess(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping concurrent benchmarks in short mode")
	}
	
	registry, tempDir := setupTestRegistry(&testing.T{})
	defer func() {
		registry.Close()
		os.RemoveAll(tempDir)
	}()
	
	// Pre-populate with test data
	numArchives := 1000
	for i := 0; i < numArchives; i++ {
		archive := &Archive{
			UID:     generateUID(),
			Name:    fmt.Sprintf("concurrent-test-%04d.7z", i),
			Path:    fmt.Sprintf("/concurrent/test-%04d.7z", i),
			Size:    int64(1024 * (i + 1)),
			Created: time.Now(),
			Profile: "balanced",
			Managed: true,
			Status:  "present",
		}
		
		err := registry.Add(archive)
		if err != nil {
			b.Fatalf("Failed to populate concurrent test data: %v", err)
		}
	}
	
	b.Run("ConcurrentReads", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			counter := 0
			for pb.Next() {
				name := fmt.Sprintf("concurrent-test-%04d.7z", counter%numArchives)
				_, err := registry.Get(name)
				if err != nil {
					b.Fatalf("Concurrent read failed: %v", err)
				}
				counter++
			}
		})
	})
	
	b.Run("ConcurrentLists", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := registry.List()
				if err != nil {
					b.Fatalf("Concurrent list failed: %v", err)
				}
			}
		})
	})
}

// BenchmarkMemoryUsage provides insights into memory consumption
func BenchmarkMemoryUsage(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping memory benchmarks in short mode")
	}
	
	registry, tempDir := setupTestRegistry(&testing.T{})
	defer func() {
		registry.Close()
		os.RemoveAll(tempDir)
	}()
	
	b.Run("MemoryScaling", func(b *testing.B) {
		// Add archives and measure memory impact
		numArchives := 10000
		
		b.ReportAllocs()
		b.ResetTimer()
		
		for i := 0; i < b.N && i < numArchives; i++ {
			archive := &Archive{
				UID:     generateUID(),
				Name:    fmt.Sprintf("memory-test-%06d.7z", i),
				Path:    fmt.Sprintf("/memory/test-%06d.7z", i),
				Size:    int64(1024 * (i + 1)),
				Created: time.Now(),
				Profile: "balanced",
				Managed: true,
				Status:  "present",
			}
			
			err := registry.Add(archive)
			if err != nil {
				b.Fatalf("Memory test add failed: %v", err)
			}
		}
	})
	
	b.Run("ListMemoryUsage", func(b *testing.B) {
		// Pre-populate with large dataset
		numArchives := 5000
		for i := 0; i < numArchives; i++ {
			archive := &Archive{
				UID:     generateUID(),
				Name:    fmt.Sprintf("list-memory-%05d.7z", i),
				Path:    fmt.Sprintf("/list/memory-%05d.7z", i),
				Size:    int64(1024 * (i + 1)),
				Created: time.Now(),
				Profile: "balanced",
				Managed: true,
				Status:  "present",
			}
			
			registry.Add(archive)
		}
		
		b.ReportAllocs()
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			archives, err := registry.List()
			if err != nil {
				b.Fatalf("List memory test failed: %v", err)
			}
			
			// Process archives to prevent optimization
			_ = len(archives)
		}
	})
}