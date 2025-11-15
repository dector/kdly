package kdly

import (
	"testing"
)

// TestParsePackageConfig tests a package manager configuration similar to Cargo.toml
func TestParsePackageConfig(t *testing.T) {
	input := `package {
	name "my-awesome-app"
	version "1.2.3"
	authors "John Doe <john@example.com>" "Jane Smith <jane@example.com>"
	license "MIT"
	description "A comprehensive example of a package configuration"
	repository url="https://github.com/example/my-awesome-app"
	keywords "kdl" "parser" "config"

	// Development dependencies
	dependencies {
		kdl version="2.0.0" features="serde" features="validation"
		serde version="1.0" optional=#true
		tokio version="1.28" features="full"
	}

	dev-dependencies {
		test-framework version="0.5.0"
		/-criterion version="0.4" // Temporarily disabled
	}

	build {
		target "x86_64-unknown-linux-gnu"
		target "aarch64-apple-darwin"
		release-flags "--opt-level" "3" "--lto"
	}
}

metadata {
	build-timestamp (date)"2025-11-15T10:30:00Z"
	commit-hash "abc123def456"
	ci-build #true
}`

	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 2 {
		t.Fatalf("Expected 2 top-level nodes, got %d", len(doc.Nodes))
	}

	// Verify package node
	pkg := doc.Nodes[0]
	if pkg.Name != "package" {
		t.Errorf("Expected 'package', got %q", pkg.Name)
	}

	if !pkg.HasChildren() {
		t.Fatal("Expected package node to have children")
	}

	// Check package name
	var nameNode *Node
	for _, child := range pkg.Children {
		if child.Name == "name" {
			nameNode = child
			break
		}
	}
	if nameNode == nil {
		t.Fatal("Expected 'name' node in package")
	}
	if len(nameNode.Arguments) != 1 || nameNode.Arguments[0].Value != "my-awesome-app" {
		t.Errorf("Expected package name 'my-awesome-app'")
	}

	// Check authors (multiple arguments)
	var authorsNode *Node
	for _, child := range pkg.Children {
		if child.Name == "authors" {
			authorsNode = child
			break
		}
	}
	if authorsNode == nil {
		t.Fatal("Expected 'authors' node in package")
	}
	if len(authorsNode.Arguments) != 2 {
		t.Errorf("Expected 2 authors, got %d", len(authorsNode.Arguments))
	}

	// Verify dependencies node
	var depsNode *Node
	for _, child := range pkg.Children {
		if child.Name == "dependencies" {
			depsNode = child
			break
		}
	}
	if depsNode == nil {
		t.Fatal("Expected 'dependencies' node")
	}

	// Check that kdl dependency has multiple features
	var kdlDep *Node
	for _, dep := range depsNode.Children {
		if dep.Name == "kdl" {
			kdlDep = dep
			break
		}
	}
	if kdlDep == nil {
		t.Fatal("Expected 'kdl' dependency")
	}
	if kdlDep.Properties["version"].Value != "2.0.0" {
		t.Errorf("Expected kdl version '2.0.0', got %v", kdlDep.Properties["version"].Value)
	}

	// Verify metadata with type annotation
	metadata := doc.Nodes[1]
	if metadata.Name != "metadata" {
		t.Errorf("Expected 'metadata', got %q", metadata.Name)
	}

	var timestampNode *Node
	for _, child := range metadata.Children {
		if child.Name == "build-timestamp" {
			timestampNode = child
			break
		}
	}
	if timestampNode == nil {
		t.Fatal("Expected 'build-timestamp' node")
	}
	if timestampNode.Arguments[0].TypeAnnotation == nil || *timestampNode.Arguments[0].TypeAnnotation != "date" {
		t.Error("Expected 'date' type annotation on timestamp")
	}
}

// TestParseCICDPipeline tests a CI/CD configuration similar to GitHub Actions
func TestParseCICDPipeline(t *testing.T) {
	input := `pipeline name="Build and Test" {
	on {
		push branches="main" branches="develop"
		pull_request
		schedule cron="0 0 * * 0" // Weekly on Sunday
	}

	env {
		NODE_ENV "production"
		CACHE_DIR "/tmp/cache"
		DEBUG #false
	}

	jobs {
		build {
			runs-on "ubuntu-latest"
			timeout-minutes 30

			steps {
				- name="Checkout code" uses="actions/checkout@v3"

				- name="Setup Node.js" uses="actions/setup-node@v3" {
					with {
						node-version "18"
						cache "npm"
					}
				}

				- name="Install dependencies" run="npm ci"

				- name="Run linter" run="npm run lint" \
					continue-on-error=#false

				- name="Build project" {
					run "npm run build"
					env {
						NODE_OPTIONS "--max-old-space-size=4096"
					}
				}

				- name="Upload artifacts" uses="actions/upload-artifact@v3" {
					with {
						name "build-output"
						path "dist/"
						retention-days 7
					}
				}
			}
		}

		test matrix-os="ubuntu-latest" matrix-os="windows-latest" {
			needs "build"
			runs-on "${{ matrix.os }}"

			steps {
				- name="Download artifacts" uses="actions/download-artifact@v3"
				- name="Run tests" run="npm test" env-CI=#true
				- name="Upload coverage" run="bash <(curl -s https://codecov.io/bash)" \
					if="success()"
			}
		}

		/-deploy {
			// Deploy job temporarily disabled
			needs "test"
			runs-on "ubuntu-latest"
		}
	}
}`

	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("Expected 1 top-level node, got %d", len(doc.Nodes))
	}

	pipeline := doc.Nodes[0]
	if pipeline.Name != "pipeline" {
		t.Errorf("Expected 'pipeline', got %q", pipeline.Name)
	}

	if pipeline.Properties["name"].Value != "Build and Test" {
		t.Errorf("Expected pipeline name 'Build and Test', got %v", pipeline.Properties["name"].Value)
	}

	// Find jobs node
	var jobsNode *Node
	for _, child := range pipeline.Children {
		if child.Name == "jobs" {
			jobsNode = child
			break
		}
	}
	if jobsNode == nil {
		t.Fatal("Expected 'jobs' node")
	}

	// Verify build job exists
	var buildJob *Node
	for _, job := range jobsNode.Children {
		if job.Name == "build" {
			buildJob = job
			break
		}
	}
	if buildJob == nil {
		t.Fatal("Expected 'build' job")
	}

	// Verify steps in build job
	var stepsNode *Node
	for _, child := range buildJob.Children {
		if child.Name == "steps" {
			stepsNode = child
			break
		}
	}
	if stepsNode == nil {
		t.Fatal("Expected 'steps' node in build job")
	}

	if len(stepsNode.Children) < 5 {
		t.Errorf("Expected at least 5 steps, got %d", len(stepsNode.Children))
	}

	// Verify that deploy job was commented out with slashdash
	deployJobCount := 0
	for _, job := range jobsNode.Children {
		if job.Name == "deploy" {
			deployJobCount++
		}
	}
	if deployJobCount != 0 {
		t.Errorf("Expected deploy job to be commented out, but found %d deploy jobs", deployJobCount)
	}

	// Verify test job has matrix properties
	var testJob *Node
	for _, job := range jobsNode.Children {
		if job.Name == "test" {
			testJob = job
			break
		}
	}
	if testJob == nil {
		t.Fatal("Expected 'test' job")
	}

	if testJob.Properties["needs"].Value != "build" {
		t.Errorf("Expected test job to need 'build', got %v", testJob.Properties["needs"].Value)
	}
}

// TestParseApplicationConfig tests a complex application configuration
func TestParseApplicationConfig(t *testing.T) {
	input := `app {
	name "production-api"
	version "3.2.1"
	debug #false

	server {
		host "0.0.0.0"
		port 8080
		read-timeout (duration)30
		write-timeout (duration)30
		max-header-bytes 1048576

		tls enabled=#true {
			cert-file "/etc/ssl/certs/server.crt"
			key-file "/etc/ssl/private/server.key"
			min-version "1.2"
			ciphers "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256" \
				"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
		}

		cors {
			allowed-origins "https://example.com" "https://app.example.com"
			allowed-methods "GET" "POST" "PUT" "DELETE" "OPTIONS"
			allowed-headers "*"
			expose-headers "X-Request-ID" "X-Response-Time"
			max-age 3600
			credentials #true
		}
	}

	database {
		primary {
			driver "postgres"
			host "db.example.com"
			port 5432
			database "production"
			username "app_user"
			password (secret)"${DB_PASSWORD}"
			ssl-mode "require"

			pool {
				max-open-conns 25
				max-idle-conns 5
				conn-max-lifetime (duration)300
				conn-max-idle-time (duration)60
			}
		}

		replica {
			driver "postgres"
			host "db-replica.example.com"
			port 5432
			database "production"
			username "readonly_user"
			password (secret)"${DB_REPLICA_PASSWORD}"
			ssl-mode "require"
			read-only #true
		}

		cache {
			driver "redis"
			addresses "redis-1.example.com:6379" \
				"redis-2.example.com:6379" \
				"redis-3.example.com:6379"
			password (secret)"${REDIS_PASSWORD}"
			db 0
			pool-size 10

			cluster-mode #true
			read-timeout (duration)3
			write-timeout (duration)3
		}
	}

	logging {
		level "info"
		format "json"
		output "stdout"

		fields {
			service "production-api"
			version "3.2.1"
			environment "production"
		}

		/*
		File logging disabled in production
		file {
			path "/var/log/app.log"
			max-size 100
			max-backups 3
			max-age 28
		}
		*/
	}

	features {
		new-dashboard enabled=#true rollout-percentage=100
		ai-suggestions enabled=#true rollout-percentage=50
		beta-api enabled=#false
		/-experimental-search enabled=#true // Disabled for now
	}

	integrations {
		email {
			provider "sendgrid"
			api-key (secret)"${SENDGRID_API_KEY}"
			from-address "noreply@example.com"
			from-name "Production API"

			templates {
				welcome template-id="d-abc123"
				reset-password template-id="d-def456"
				notification template-id="d-ghi789"
			}
		}

		monitoring {
			sentry {
				dsn (url)"https://xxx@yyy.ingest.sentry.io/zzz"
				environment "production"
				traces-sample-rate (f64)0.1
				profiles-sample-rate (f64)0.1
			}

			prometheus {
				enabled #true
				path "/metrics"
				namespace "production_api"
			}
		}

		storage {
			s3 {
				region "us-east-1"
				bucket "production-uploads"
				access-key-id (secret)"${AWS_ACCESS_KEY_ID}"
				secret-access-key (secret)"${AWS_SECRET_ACCESS_KEY}"

				encryption {
					type "AES256"
					kms-key-id "arn:aws:kms:us-east-1:123456789:key/abc-def"
				}
			}
		}
	}

	rate-limiting {
		enabled #true

		tiers {
			free requests-per-minute=60 burst=10
			pro requests-per-minute=600 burst=100
			enterprise requests-per-minute=6000 burst=1000
		}

		storage "redis"
		key-prefix "ratelimit:"
	}
}`

	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("Expected 1 top-level node, got %d", len(doc.Nodes))
	}

	app := doc.Nodes[0]
	if app.Name != "app" {
		t.Errorf("Expected 'app', got %q", app.Name)
	}

	// Verify server configuration
	var serverNode *Node
	for _, child := range app.Children {
		if child.Name == "server" {
			serverNode = child
			break
		}
	}
	if serverNode == nil {
		t.Fatal("Expected 'server' node")
	}

	// Check server port
	var portNode *Node
	for _, child := range serverNode.Children {
		if child.Name == "port" {
			portNode = child
			break
		}
	}
	if portNode == nil {
		t.Fatal("Expected 'port' node")
	}
	if len(portNode.Arguments) != 1 || portNode.Arguments[0].Value != int64(8080) {
		t.Errorf("Expected port 8080")
	}

	// Verify TLS configuration
	var tlsNode *Node
	for _, child := range serverNode.Children {
		if child.Name == "tls" {
			tlsNode = child
			break
		}
	}
	if tlsNode == nil {
		t.Fatal("Expected 'tls' node")
	}
	if tlsNode.Properties["enabled"].Value != true {
		t.Error("Expected TLS to be enabled")
	}

	// Verify database configuration with multiple sections
	var dbNode *Node
	for _, child := range app.Children {
		if child.Name == "database" {
			dbNode = child
			break
		}
	}
	if dbNode == nil {
		t.Fatal("Expected 'database' node")
	}

	// Should have primary, replica, and cache
	if len(dbNode.Children) != 3 {
		t.Errorf("Expected 3 database configurations, got %d", len(dbNode.Children))
	}

	// Verify cache has type annotation for timeouts
	var cacheNode *Node
	for _, child := range dbNode.Children {
		if child.Name == "cache" {
			cacheNode = child
			break
		}
	}
	if cacheNode == nil {
		t.Fatal("Expected 'cache' node")
	}

	var readTimeoutNode *Node
	for _, child := range cacheNode.Children {
		if child.Name == "read-timeout" {
			readTimeoutNode = child
			break
		}
	}
	if readTimeoutNode == nil {
		t.Fatal("Expected 'read-timeout' node in cache")
	}
	if readTimeoutNode.Arguments[0].TypeAnnotation == nil || *readTimeoutNode.Arguments[0].TypeAnnotation != "duration" {
		t.Error("Expected 'duration' type annotation on read-timeout")
	}

	// Verify features with rollout percentages
	var featuresNode *Node
	for _, child := range app.Children {
		if child.Name == "features" {
			featuresNode = child
			break
		}
	}
	if featuresNode == nil {
		t.Fatal("Expected 'features' node")
	}

	// Check that experimental-search was commented out
	experimentalSearchCount := 0
	for _, feature := range featuresNode.Children {
		if feature.Name == "experimental-search" {
			experimentalSearchCount++
		}
	}
	if experimentalSearchCount != 0 {
		t.Errorf("Expected experimental-search to be commented out, but found %d instances", experimentalSearchCount)
	}

	// Verify integrations exist
	var integrationsNode *Node
	for _, child := range app.Children {
		if child.Name == "integrations" {
			integrationsNode = child
			break
		}
	}
	if integrationsNode == nil {
		t.Fatal("Expected 'integrations' node")
	}

	// Should have email, monitoring, and storage
	if len(integrationsNode.Children) != 3 {
		t.Errorf("Expected 3 integration sections, got %d", len(integrationsNode.Children))
	}

	// Verify monitoring has nested sentry and prometheus
	var monitoringNode *Node
	for _, child := range integrationsNode.Children {
		if child.Name == "monitoring" {
			monitoringNode = child
			break
		}
	}
	if monitoringNode == nil {
		t.Fatal("Expected 'monitoring' node")
	}
	if len(monitoringNode.Children) != 2 {
		t.Errorf("Expected 2 monitoring integrations, got %d", len(monitoringNode.Children))
	}

	// Verify type annotations on sensitive data
	var sentryNode *Node
	for _, child := range monitoringNode.Children {
		if child.Name == "sentry" {
			sentryNode = child
			break
		}
	}
	if sentryNode == nil {
		t.Fatal("Expected 'sentry' node")
	}

	var tracesNode *Node
	for _, child := range sentryNode.Children {
		if child.Name == "traces-sample-rate" {
			tracesNode = child
			break
		}
	}
	if tracesNode == nil {
		t.Fatal("Expected 'traces-sample-rate' node")
	}
	if tracesNode.Arguments[0].TypeAnnotation == nil || *tracesNode.Arguments[0].TypeAnnotation != "f64" {
		t.Error("Expected 'f64' type annotation on traces-sample-rate")
	}

	// Verify rate-limiting tiers
	var rateLimitNode *Node
	for _, child := range app.Children {
		if child.Name == "rate-limiting" {
			rateLimitNode = child
			break
		}
	}
	if rateLimitNode == nil {
		t.Fatal("Expected 'rate-limiting' node")
	}

	var tiersNode *Node
	for _, child := range rateLimitNode.Children {
		if child.Name == "tiers" {
			tiersNode = child
			break
		}
	}
	if tiersNode == nil {
		t.Fatal("Expected 'tiers' node")
	}
	if len(tiersNode.Children) != 3 {
		t.Errorf("Expected 3 rate limit tiers, got %d", len(tiersNode.Children))
	}
}

// TestParseWebServerRoutes tests a web server routing configuration
func TestParseWebServerRoutes(t *testing.T) {
	input := `routes {
	group path="/api/v1" {
		middleware "auth" "logging" "cors"

		GET path="/users" handler="listUsers" {
			query-params {
				page type="int" default=1
				limit type="int" default=20 max=100
				sort type="string" enum="name" enum="created" enum="updated"
			}
			response {
				200 type="application/json" schema="UserList"
				401 type="application/json" schema="Error"
			}
		}

		POST path="/users" handler="createUser" {
			middleware "rate-limit" config-key="user-creation"

			request {
				type "application/json"
				schema "CreateUserRequest"
				required #true
			}

			response {
				201 type="application/json" schema="User"
				400 type="application/json" schema="ValidationError"
				401 type="application/json" schema="Error"
			}
		}

		GET path="/users/:id" handler="getUser" {
			path-params {
				id type="uuid" required=#true
			}

			response {
				200 type="application/json" schema="User"
				404 type="application/json" schema="Error"
			}
		}

		PUT path="/users/:id" handler="updateUser" {
			middleware "rate-limit"

			path-params {
				id type="uuid" required=#true
			}

			request {
				type "application/json"
				schema "UpdateUserRequest"
			}

			response {
				200 type="application/json" schema="User"
				404 type="application/json" schema="Error"
			}
		}

		DELETE path="/users/:id" handler="deleteUser" {
			middleware "admin-only"

			response {
				204 type #null
				403 type="application/json" schema="Error"
				404 type="application/json" schema="Error"
			}
		}

		// File upload endpoint
		POST path="/upload" handler="handleUpload" {
			middleware "auth" "file-size-limit"

			request {
				type "multipart/form-data"
				max-size (bytes)10485760 // 10MB
			}

			response {
				200 type="application/json" schema="UploadResponse"
				413 type="application/json" schema="Error"
			}
		}
	}

	group path="/admin" {
		middleware "auth" "admin-only" "audit-log"

		GET path="/stats" handler="getStats"
		POST path="/users/:id/ban" handler="banUser"

		/-GET path="/debug" handler="debugInfo" // Disabled in production
	}

	// Static file serving
	static {
		path "/assets"
		root "./public/assets"
		index "index.html"
		spa-mode #true

		cache {
			enabled #true
			max-age 86400
			patterns "*.js" "*.css" "*.woff2"
		}
	}

	/* Legacy API endpoints - to be removed in v2
	group path="/api/v0" {
		GET path="/legacy" handler="legacyHandler"
	}
	*/
}`

	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("Expected 1 top-level node, got %d", len(doc.Nodes))
	}

	routes := doc.Nodes[0]
	if routes.Name != "routes" {
		t.Errorf("Expected 'routes', got %q", routes.Name)
	}

	// Should have 2 group nodes and 1 static node (legacy group is commented out)
	if len(routes.Children) != 3 {
		t.Errorf("Expected 3 child nodes, got %d", len(routes.Children))
	}

	// Verify first group (API v1)
	apiGroup := routes.Children[0]
	if apiGroup.Name != "group" {
		t.Errorf("Expected 'group', got %q", apiGroup.Name)
	}
	if apiGroup.Properties["path"].Value != "/api/v1" {
		t.Errorf("Expected path '/api/v1', got %v", apiGroup.Properties["path"].Value)
	}

	// Count HTTP method nodes in API group (should have GET, POST, GET, PUT, DELETE, POST)
	methodCount := 0
	for _, child := range apiGroup.Children {
		if child.Name == "GET" || child.Name == "POST" || child.Name == "PUT" || child.Name == "DELETE" {
			methodCount++
		}
	}
	if methodCount != 6 {
		t.Errorf("Expected 6 HTTP method nodes in API group, got %d", methodCount)
	}

	// Verify admin group
	adminGroup := routes.Children[1]
	if adminGroup.Name != "group" {
		t.Errorf("Expected 'group', got %q", adminGroup.Name)
	}
	if adminGroup.Properties["path"].Value != "/admin" {
		t.Errorf("Expected path '/admin', got %v", adminGroup.Properties["path"].Value)
	}

	// Verify debug endpoint was commented out
	debugCount := 0
	for _, child := range adminGroup.Children {
		if child.Name == "GET" && child.Properties["path"] != nil && child.Properties["path"].Value == "/debug" {
			debugCount++
		}
	}
	if debugCount != 0 {
		t.Errorf("Expected debug endpoint to be commented out, but found %d instances", debugCount)
	}

	// Verify static configuration
	staticNode := routes.Children[2]
	if staticNode.Name != "static" {
		t.Errorf("Expected 'static', got %q", staticNode.Name)
	}

	// Check spa-mode property
	var spaModeNode *Node
	for _, child := range staticNode.Children {
		if child.Name == "spa-mode" {
			spaModeNode = child
			break
		}
	}
	if spaModeNode == nil {
		t.Fatal("Expected 'spa-mode' node")
	}
	if len(spaModeNode.Arguments) != 1 || spaModeNode.Arguments[0].Value != true {
		t.Error("Expected spa-mode to be true")
	}

	// Verify cache configuration exists
	var cacheNode *Node
	for _, child := range staticNode.Children {
		if child.Name == "cache" {
			cacheNode = child
			break
		}
	}
	if cacheNode == nil {
		t.Fatal("Expected 'cache' node in static")
	}
	if !cacheNode.HasChildren() {
		t.Error("Expected cache node to have children")
	}
}

// TestParseBuildSystemConfig tests a complex build system configuration
func TestParseBuildSystemConfig(t *testing.T) {
	input := `build {
	project "multiplatform-app"
	version "2.1.0"

	toolchain {
		go version="1.21"
		node version="18.17"
		python version="3.11"
	}

	targets {
		binary name="app-server" {
			type "executable"
			main "./cmd/server"
			output "bin/server"

			platforms {
				linux arch="amd64" arch="arm64"
				darwin arch="amd64" arch="arm64"
				windows arch="amd64"
			}

			build-flags {
				ldflags "-s -w -X main.Version=${VERSION}"
				tags "netgo" "osusergo"
			}

			dependencies {
				internal "./pkg/api" "./pkg/auth" "./pkg/db"
				external "github.com/gorilla/mux" version="1.8.0"
				external "github.com/lib/pq" version="1.10.9"
			}
		}

		binary name="app-cli" {
			type "executable"
			main "./cmd/cli"
			output "bin/cli"

			platforms {
				linux arch="amd64"
				darwin arch="amd64" arch="arm64"
			}
		}

		library name="app-sdk" {
			type "shared"
			main "./sdk"
			output "lib/libapp"

			exports {
				Initialize
				Configure
				Execute
				Shutdown
			}

			platforms {
				linux arch="amd64" extension=".so"
				darwin arch="arm64" extension=".dylib"
				windows arch="amd64" extension=".dll"
			}
		}

		container name="app-docker" {
			type "docker"
			dockerfile "./Dockerfile"
			context "."

			tags {
				tag "myapp:latest"
				tag "myapp:${VERSION}"
				tag "myapp:${GIT_SHA}"
			}

			build-args {
				GO_VERSION "1.21"
				NODE_VERSION "18"
				BUILD_DATE "${BUILD_DATE}"
			}

			platforms {
				linux-amd64 platform="linux/amd64"
				linux-arm64 platform="linux/arm64"
			}
		}
	}

	tasks {
		clean {
			description "Remove build artifacts"
			commands {
				- "rm -rf bin/"
				- "rm -rf lib/"
				- "rm -rf dist/"
			}
		}

		test {
			description "Run all tests"
			depends-on "lint"

			commands {
				- "go test -v -race -coverprofile=coverage.txt ./..."
				- "go tool cover -html=coverage.txt -o coverage.html"
			}

			env {
				GO111MODULE "on"
				CGO_ENABLED "1"
			}
		}

		lint {
			description "Run linters"

			commands {
				- "golangci-lint run ./..."
				- "npm run lint"
			}
		}

		build {
			description "Build all targets"
			depends-on "test"
			parallel #true

			steps {
				step target="app-server" platforms="linux/amd64" platforms="darwin/arm64"
				step target="app-cli" platforms="all"
				step target="app-sdk" platforms="linux/amd64"
			}
		}

		release {
			description "Create release artifacts"
			depends-on "build"

			commands {
				- name="Package binaries" run="./scripts/package.sh"
				- name="Generate checksums" run="sha256sum dist/* > dist/checksums.txt"
				- name="Sign artifacts" run="./scripts/sign.sh" if="${SIGN_ARTIFACTS}"
			}

			artifacts {
				- path="dist/*.tar.gz"
				- path="dist/*.zip"
				- path="dist/checksums.txt"
			}
		}

		/-benchmark {
			// Benchmarks disabled for regular builds
			description "Run performance benchmarks"
			commands {
				- "go test -bench=. -benchmem ./..."
			}
		}
	}

	environments {
		development {
			DEBUG #true
			LOG_LEVEL "debug"
			API_ENDPOINT "http://localhost:8080"
		}

		staging {
			DEBUG #false
			LOG_LEVEL "info"
			API_ENDPOINT "https://staging-api.example.com"
			RATE_LIMIT 100
		}

		production {
			DEBUG #false
			LOG_LEVEL "warn"
			API_ENDPOINT "https://api.example.com"
			RATE_LIMIT 1000
			ENABLE_METRICS #true
			ENABLE_TRACING #true
		}
	}

	/*
	Future improvements:
	- Add caching for build artifacts
	- Implement incremental builds
	- Add support for cross-compilation caching
	*/
}`

	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("Expected 1 top-level node, got %d", len(doc.Nodes))
	}

	build := doc.Nodes[0]
	if build.Name != "build" {
		t.Errorf("Expected 'build', got %q", build.Name)
	}

	// Verify targets section
	var targetsNode *Node
	for _, child := range build.Children {
		if child.Name == "targets" {
			targetsNode = child
			break
		}
	}
	if targetsNode == nil {
		t.Fatal("Expected 'targets' node")
	}

	// Should have 3 targets: 2 binaries and 1 library and 1 container
	if len(targetsNode.Children) != 4 {
		t.Errorf("Expected 4 targets, got %d", len(targetsNode.Children))
	}

	// Verify app-server binary
	var serverBinary *Node
	for _, target := range targetsNode.Children {
		if target.Name == "binary" && target.Properties["name"] != nil && target.Properties["name"].Value == "app-server" {
			serverBinary = target
			break
		}
	}
	if serverBinary == nil {
		t.Fatal("Expected 'app-server' binary target")
	}

	// Check platforms
	var platformsNode *Node
	for _, child := range serverBinary.Children {
		if child.Name == "platforms" {
			platformsNode = child
			break
		}
	}
	if platformsNode == nil {
		t.Fatal("Expected 'platforms' node in app-server")
	}
	if len(platformsNode.Children) != 3 {
		t.Errorf("Expected 3 platforms for app-server, got %d", len(platformsNode.Children))
	}

	// Verify tasks section
	var tasksNode *Node
	for _, child := range build.Children {
		if child.Name == "tasks" {
			tasksNode = child
			break
		}
	}
	if tasksNode == nil {
		t.Fatal("Expected 'tasks' node")
	}

	// Verify benchmark task was commented out
	benchmarkCount := 0
	for _, task := range tasksNode.Children {
		if task.Name == "benchmark" {
			benchmarkCount++
		}
	}
	if benchmarkCount != 0 {
		t.Errorf("Expected benchmark task to be commented out, but found %d instances", benchmarkCount)
	}

	// Verify build task has parallel flag
	var buildTask *Node
	for _, task := range tasksNode.Children {
		if task.Name == "build" {
			buildTask = task
			break
		}
	}
	if buildTask == nil {
		t.Fatal("Expected 'build' task")
	}
	if buildTask.Properties["parallel"] == nil || buildTask.Properties["parallel"].Value != true {
		t.Error("Expected build task to have parallel=true")
	}

	// Verify environments section
	var envsNode *Node
	for _, child := range build.Children {
		if child.Name == "environments" {
			envsNode = child
			break
		}
	}
	if envsNode == nil {
		t.Fatal("Expected 'environments' node")
	}

	// Should have 3 environments
	if len(envsNode.Children) != 3 {
		t.Errorf("Expected 3 environments, got %d", len(envsNode.Children))
	}

	// Verify production environment
	var prodEnv *Node
	for _, env := range envsNode.Children {
		if env.Name == "production" {
			prodEnv = env
			break
		}
	}
	if prodEnv == nil {
		t.Fatal("Expected 'production' environment")
	}

	var metricsNode *Node
	for _, child := range prodEnv.Children {
		if child.Name == "ENABLE_METRICS" {
			metricsNode = child
			break
		}
	}
	if metricsNode == nil {
		t.Fatal("Expected 'ENABLE_METRICS' node in production environment")
	}
	if len(metricsNode.Arguments) != 1 || metricsNode.Arguments[0].Value != true {
		t.Error("Expected ENABLE_METRICS to be true")
	}
}
