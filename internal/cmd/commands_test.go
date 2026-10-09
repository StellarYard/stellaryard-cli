package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// commandNames returns the names of every subcommand registered on rootCmd.
func commandNames() map[string]bool {
	names := map[string]bool{}
	for _, c := range rootCmd.Commands() {
		names[c.Name()] = true
	}
	return names
}

func TestRootRegistersAllDocumentedCommands(t *testing.T) {
	names := commandNames()
	for _, want := range []string{"containers", "accounts", "ledger", "contracts", "logs"} {
		if !names[want] {
			t.Errorf("root command is missing the %q subcommand; registered: %v", want, names)
		}
	}
}

func TestAccountsExposesOnlyDocumentedSubcommands(t *testing.T) {
	// The README documents `accounts create` and `accounts list`. If a
	// subcommand is added or removed, update the README in the same commit.
	var found []string
	for _, c := range accountsCmd.Commands() {
		found = append(found, c.Name())
	}
	for _, want := range []string{"create", "list"} {
		if !contains(found, want) {
			t.Errorf("accounts is missing %q; has %v", want, found)
		}
	}
	if contains(found, "show") {
		t.Error("accounts show exists but is not documented as implemented — update the README")
	}
}

func TestExitCodesMatchReadmeTable(t *testing.T) {
	// README documents: 0 success, 1 argument error, 2 core unreachable,
	// 3 core app error.
	cases := map[string]int{
		"ExitSuccess":         ExitSuccess,
		"ExitArgError":        ExitArgError,
		"ExitCoreUnreachable": ExitCoreUnreachable,
		"ExitAppError":        ExitAppError,
	}
	want := map[string]int{
		"ExitSuccess":         0,
		"ExitArgError":        1,
		"ExitCoreUnreachable": 2,
		"ExitAppError":        3,
	}
	for name, got := range cases {
		if got != want[name] {
			t.Errorf("%s = %d, README documents %d", name, got, want[name])
		}
	}
}

func TestPersistentFlagsDefaults(t *testing.T) {
	cases := []struct {
		name  string
		want  string
		deflt string
	}{
		{"core-url", "http://localhost:8080", "http://localhost:8080"},
		{"format", "table", "table"},
		{"api-key", "", ""},
	}
	for _, c := range cases {
		f := rootCmd.PersistentFlags().Lookup(c.name)
		if f == nil {
			t.Errorf("persistent flag --%s is not registered", c.name)
			continue
		}
		if f.DefValue != c.deflt {
			t.Errorf("--%s default = %q, want %q", c.name, f.DefValue, c.deflt)
		}
	}
	if f := rootCmd.PersistentFlags().Lookup("verbose"); f == nil {
		t.Error("persistent flag --verbose is not registered")
	}
}

func TestAPIKeyAuthorizationHeaderForwarding(t *testing.T) {
	var authHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `[]`)
	}))
	defer srv.Close()

	oldURL, oldKey := coreURL, apiKey
	coreURL = srv.URL
	apiKey = "test-secret-key-123"
	t.Cleanup(func() {
		coreURL, apiKey = oldURL, oldKey
	})

	captureStdout(t, func() {
		containersStatusCmd.Run(containersStatusCmd, nil)
	})

	if authHeader != "Bearer test-secret-key-123" {
		t.Errorf("Authorization header = %q, want %q", authHeader, "Bearer test-secret-key-123")
	}
}

func TestContainerCommandsRequireNameFlag(t *testing.T) {
	start := containersStartCmd.Flags().Lookup("name")
	if start == nil {
		t.Error("containers start is missing --name")
	} else if start.DefValue != "" {
		t.Errorf("containers start --name default = %q, want empty", start.DefValue)
	}

	stop := containersStopCmd.Flags().Lookup("name")
	if stop == nil {
		t.Error("containers stop is missing --name")
	} else if stop.DefValue != "" {
		t.Errorf("containers stop --name default = %q, want empty", stop.DefValue)
	}

	// Both commands use --name rather than a positional argument, because
	// the README documents `containers start --name horizon`.
	if containersStartCmd.Args != nil || containersStopCmd.Args != nil {
		t.Error("container start/stop should take --name, not positional args")
	}
}

func TestLogsFollowFlagExists(t *testing.T) {
	f := logsCmd.Flags().Lookup("follow")
	if f == nil {
		t.Fatal("logs is missing --follow")
	}
	if f.DefValue != "false" {
		t.Errorf("--follow default = %q, want false", f.DefValue)
	}
}

func TestAccountsCreateLabelFlagExists(t *testing.T) {
	f := accountsCreateCmd.Flags().Lookup("label")
	if f == nil {
		t.Fatal("accounts create is missing --label")
	}
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// captureStdout redirects os.Stdout for the duration of fn and returns what
// was written.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read captured output: %v", err)
	}
	return string(out)
}

// withGlobals sets the package-level command state and restores it afterwards.
func withGlobals(t *testing.T, url, format string) {
	t.Helper()
	oldURL, oldFormat := coreURL, formatStr
	coreURL, formatStr = url, format
	t.Cleanup(func() { coreURL, formatStr = oldURL, oldFormat })
}

func TestContainersStatusTableOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/containers" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `[{"name":"horizon","state":"running"},{"name":"soroban-rpc","state":"exited"}]`)
	}))
	defer srv.Close()
	withGlobals(t, srv.URL, "table")

	out := captureStdout(t, func() {
		containersStatusCmd.Run(containersStatusCmd, nil)
	})

	for _, want := range []string{"NAME", "horizon", "running", "soroban-rpc", "exited"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}

func TestContainersStatusJSONOutput(t *testing.T) {
	payload := `[{"name":"horizon","state":"running"}]`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, payload)
	}))
	defer srv.Close()
	withGlobals(t, srv.URL, "json")

	out := captureStdout(t, func() {
		containersStatusCmd.Run(containersStatusCmd, nil)
	})

	var got []map[string]string
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &got); err != nil {
		t.Fatalf("--format json output is not valid JSON: %v\n%s", err, out)
	}
	if len(got) != 1 || got[0]["name"] != "horizon" {
		t.Errorf("parsed %v, want one horizon container", got)
	}
}

func TestContainersStatusEmptyList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `[]`)
	}))
	defer srv.Close()
	withGlobals(t, srv.URL, "table")

	out := captureStdout(t, func() {
		containersStatusCmd.Run(containersStatusCmd, nil)
	})

	if !strings.Contains(out, "No containers found") {
		t.Errorf("empty list should say 'No containers found', got:\n%s", out)
	}
}

func TestAccountsListJSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		io.WriteString(w, `[{"id":"1","publicKey":"GABC","label":"deployer","network":"local"}]`)
	}))
	defer srv.Close()
	withGlobals(t, srv.URL, "json")

	out := captureStdout(t, func() {
		accountsListCmd.Run(accountsListCmd, nil)
	})

	var got []map[string]string
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &got); err != nil {
		t.Fatalf("--format json output is not valid JSON: %v\n%s", err, out)
	}
	if len(got) != 1 || got[0]["publicKey"] != "GABC" {
		t.Errorf("parsed %v, want one account with publicKey GABC", got)
	}
}

func TestAccountsCreateSendsLabel(t *testing.T) {
	var receivedLabel string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/accounts" {
			t.Errorf("got %s %s, want POST /api/v1/accounts", r.Method, r.URL.Path)
		}
		var body struct {
			Label string `json:"label"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("request body is not valid JSON: %v", err)
		}
		receivedLabel = body.Label
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"id":"1","publicKey":"GABC","label":"`+body.Label+`","network":"local"}`)
	}))
	defer srv.Close()
	withGlobals(t, srv.URL, "json")

	if err := accountsCreateCmd.Flags().Set("label", "from-test"); err != nil {
		t.Fatalf("set --label: %v", err)
	}
	t.Cleanup(func() { accountsCreateCmd.Flags().Set("label", "") })

	captureStdout(t, func() {
		accountsCreateCmd.Run(accountsCreateCmd, nil)
	})

	if receivedLabel != "from-test" {
		t.Errorf("server received label %q, want %q", receivedLabel, "from-test")
	}
}

func TestLogsStubDoesNotClaimSuccess(t *testing.T) {
	// logs is a documented stub until core ships its WebSocket endpoint.
	// This guards against it silently pretending to stream.
	out := captureStdout(t, func() {
		if err := logsCmd.RunE(logsCmd, []string{"horizon"}); err != nil {
			t.Errorf("logs RunE returned error: %v", err)
		}
	})
	if !strings.Contains(out, "not implemented") {
		t.Errorf("logs output should state it is not implemented, got:\n%s", out)
	}
}
