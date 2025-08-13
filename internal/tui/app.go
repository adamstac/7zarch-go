package tui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/display"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type archivesMsg struct {
	archives []*storage.Archive
	err      error
}

type App struct {
	tbl        table.Model
	archives   []*storage.Archive
	width      int
	height     int
	statusLine string
	showHelp   bool
	showDetail bool

	// Sorting & filtering
	sortBy       string
	sortAsc      bool
	statusFilter string

	// Selection
	selected map[string]bool

	// Actions overlay
	showActions   bool
	actionIndex   int
	showConfirm   bool
	confirmAction string

	// Errors overlay
	showErrors bool
	lastErrors []string
}

func NewApp() *App {
	columns := []table.Column{
		{Title: "", Width: 2},
		{Title: "ID", Width: 12},
		{Title: "Name", Width: 30},
		{Title: "Size", Width: 10},
		{Title: "Status", Width: 8},
		{Title: "Created", Width: 19},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	// Styles
	header := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("63")).Foreground(lipgloss.Color("230"))
	selected := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)

	t.SetStyles(table.Styles{
		Header:   header,
		Cell:     lipgloss.NewStyle(),
		Selected: selected,
	})
	// Wrap the table with a base box when rendering by adjusting width/height, but Styles.Base
	// is not available in the current table version; we'll use lipgloss around View instead.

	return &App{tbl: t, sortBy: "created", sortAsc: false, statusFilter: "all", selected: map[string]bool{}}
}

func (a *App) Init() tea.Cmd { return a.loadArchivesCmd() }

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		// When confirm dialog is open
		if a.showConfirm {
			switch m.String() {
			case "enter", "y":
				// execute selected action
				a.executeAction()
				a.showConfirm = false
				a.showActions = false
				return a, a.loadArchivesCmd()
			case "esc", "n":
				a.showConfirm = false
				return a, nil
			}
			return a, nil
		}
		// When actions overlay is open
		if a.showActions {
			switch m.String() {
			case "esc", "a":
				a.showActions = false
				return a, nil
			case "up", "k":
				if a.actionIndex > 0 {
					a.actionIndex--
				}
				return a, nil
			case "down", "j":
				if a.actionIndex < 4 {
					a.actionIndex++
				}
				return a, nil
			case "enter":
				a.prepareAction()
				return a, nil
			}
			return a, nil
		}
		switch m.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		case "up", "k":
			a.tbl.MoveUp(1)
			return a, nil
		case "down", "j":
			a.tbl.MoveDown(1)
			return a, nil
		case "home", "g":
			a.tbl.GotoTop()
			return a, nil
		case "end", "G":
			a.tbl.GotoBottom()
			return a, nil
		case "a":
			// toggle actions overlay
			a.showActions = !a.showActions
			if a.showActions {
				a.actionIndex = 0
			}
			return a, nil
		case "s":
			a.cycleSortField()
			a.refreshRows()
			return a, nil
		case "o":
			a.sortAsc = !a.sortAsc
			a.refreshRows()
			return a, nil
		case "f":
			a.cycleStatusFilter()
			a.refreshRows()
			return a, nil
		case " ":
			if list := a.filtered(); len(list) > 0 && a.tbl.Cursor() >= 0 && a.tbl.Cursor() < len(list) {
				id := list[a.tbl.Cursor()].UID
				if a.selected == nil {
					a.selected = map[string]bool{}
				}
				a.selected[id] = !a.selected[id]
			}
			return a, nil
		case "r":
			return a, a.loadArchivesCmd()
		case "?":
			a.showHelp = !a.showHelp
			return a, nil
		case "e":
			if len(a.lastErrors) > 0 {
				a.showErrors = !a.showErrors
			}
			return a, nil
		case "enter":
			a.showDetail = !a.showDetail
			return a, nil
		case "esc":
			a.showDetail = false
			return a, nil
		}
	case tea.WindowSizeMsg:
		a.width, a.height = m.Width, m.Height
		// Keep some space for details/help/status
		rows := a.height - 6
		if rows < 5 {
			rows = 5
		}
		a.tbl.SetHeight(rows)
		// Set a reasonable total width; columns are fixed widths
		a.tbl.SetWidth(min(a.width-2, 90))
	case archivesMsg:
		a.statusLine = ""
		if m.err != nil {
			a.statusLine = fmt.Sprintf("Error: %v", m.err)
			return a, nil
		}
		a.archives = m.archives
		a.refreshRows()
		a.statusLine = fmt.Sprintf("Loaded %d archives", len(m.archives))
	}
	// Let table handle any remaining messages (pagination, etc.)
	var cmd tea.Cmd
	a.tbl, cmd = a.tbl.Update(msg)
	return a, cmd
}

func (a *App) View() string {
	// Title bar
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63")).Render("7zarch-go — Archives")
	// legend: "  ↑/↓ or j/k move • enter details • r refresh • ? help • q quit"

	out := lipgloss.JoinHorizontal(lipgloss.Center, title)
	out += "\n"
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(0, 1)
	out += box.Render(a.tbl.View())

	// Details pane
	if a.showDetail {
		list := a.filtered()
		if len(list) > 0 && a.tbl.Cursor() >= 0 && a.tbl.Cursor() < len(list) {
			arc := list[a.tbl.Cursor()]
			detail := a.renderDetails(arc)
			out += "\n" + detail
		}
	}

	// Footer: Confirm, Errors, Actions, Help or Status
	if a.showConfirm {
		out += "\n" + a.renderConfirm()
	} else if a.showErrors {
		out += "\n" + a.renderErrors()
	} else if a.showActions {
		out += "\n" + a.renderActions()
	} else if a.showHelp {
		help := a.renderHelp()
		out += "\n" + help
	} else {
		status := lipgloss.NewStyle().Faint(true).Render(a.statusSummary())
		out += "\n" + status
	}

	return out
}

func (a *App) rowsFromArchives(archives []*storage.Archive) []table.Row {
	rows := make([]table.Row, 0, len(archives))
	for _, arc := range archives {
		id := arc.UID
		if len(id) > 12 {
			id = id[:12]
		}
		name := arc.Name
		if len(name) > 30 {
			name = name[:29] + "…"
		}
		sel := ""
		if a.selected != nil && a.selected[arc.UID] {
			sel = "✓"
		}
		size := display.FormatSize(arc.Size)
		status := display.FormatStatus(arc.Status, false)
		created := arc.Created.Format("2006-01-02 15:04:05")
		rows = append(rows, table.Row{sel, id, name, size, status, created})
	}
	return rows
}

func (a *App) renderDetails(arc *storage.Archive) string {
	label := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render
	val := lipgloss.NewStyle().Bold(true).Render
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("99")).Padding(0, 1)

	lines := []string{
		label("Details"),
		fmt.Sprintf("ID: %s", val(arc.UID)),
		fmt.Sprintf("Name: %s", val(arc.Name)),
		fmt.Sprintf("Path: %s", val(arc.Path)),
		fmt.Sprintf("Size: %s", val(display.FormatSize(arc.Size))),
		fmt.Sprintf("Created: %s", val(arc.Created.Format("2006-01-02 15:04:05"))),
		fmt.Sprintf("Status: %s", val(display.FormatStatus(arc.Status, false))),
		fmt.Sprintf("Managed: %v", val(fmt.Sprintf("%v", arc.Managed))),
	}
	content := ""
	for _, l := range lines {
		content += l + "\n"
	}
	return box.Width(min(a.width-2, 90)).Render(content)
}

func (a *App) prepareAction() {
	items := []string{"delete", "restore", "move", "mark-uploaded", "cancel"}
	if a.actionIndex < 0 || a.actionIndex >= len(items) {
		a.showActions = false
		return
	}
	a.confirmAction = items[a.actionIndex]
	if a.confirmAction == "cancel" {
		a.showActions = false
		return
	}
	a.showConfirm = true
}

func (a *App) executeAction() {
	// For now, only implement delete/restore via existing CLI patterns using registry
	cfg, _ := config.Load()
	mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		a.statusLine = fmt.Sprintf("error: %v", err)
		return
	}
	defer mgr.Close()
	list := a.filtered()
	if len(list) == 0 {
		return
	}
	// If multiple selected, operate on those; else use cursor
	uids := []string{}
	for _, arc := range list {
		if a.selected[arc.UID] {
			uids = append(uids, arc.UID)
		}
	}
	if len(uids) == 0 {
		if a.tbl.Cursor() >= 0 && a.tbl.Cursor() < len(list) {
			uids = append(uids, list[a.tbl.Cursor()].UID)
		}
	}
	resolver := storage.NewResolver(mgr.Registry())
	a.lastErrors = nil
	count := 0
	for _, id := range uids {
		arc, err := resolver.Resolve(id)
		if err != nil {
			a.lastErrors = append(a.lastErrors, fmt.Sprintf("resolve %s: %v", id, err))
			continue
		}
		switch a.confirmAction {
		case "delete":
			// Soft delete like cmd/mas_delete.go
			now := time.Now()
			orig := arc.Path
			if arc.Managed {
				trashDir := mgr.GetTrashPath()
				_ = os.MkdirAll(trashDir, 0750)
				trashPath := filepath.Join(trashDir, filepath.Base(arc.Path))
				if err := moveOrCopy(arc.Path, trashPath); err != nil {
					a.lastErrors = append(a.lastErrors, fmt.Sprintf("move to trash %s: %v", arc.Name, err))
				} else {
					arc.Path = trashPath
				}
			}
			arc.Status = "deleted"
			arc.DeletedAt = &now
			if arc.OriginalPath == "" {
				arc.OriginalPath = orig
			}
			_ = mgr.Registry().Update(arc)
			count++
		case "restore":
			if arc.Status != "deleted" {
				continue
			}
			target := arc.OriginalPath
			if target == "" {
				name := arc.Name
				if name == "" {
					name = filepath.Base(arc.Path)
				}
				target = mgr.GetManagedPath(name)
			}
			_ = os.MkdirAll(filepath.Dir(target), 0750)
			if err := moveOrCopy(arc.Path, target); err != nil {
				a.lastErrors = append(a.lastErrors, fmt.Sprintf("restore move %s: %v", arc.Name, err))
				continue
			}
			arc.Path = target
			arc.Status = "present"
			arc.DeletedAt = nil
			_ = mgr.Registry().Update(arc)
			count++
		}
	}
	a.statusLine = fmt.Sprintf("%s: %d item(s)", a.confirmAction, count)
}

// moveOrCopy tries to rename; if it fails (e.g., cross-device), it copies then removes
func moveOrCopy(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()
	dstF, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer dstF.Close()
	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}
	return os.Remove(src)
}

func (a *App) renderConfirm() string {
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("204")).Padding(0, 1)
	msg := fmt.Sprintf("Confirm %s? (enter/y = yes, esc/n = no)", a.confirmAction)
	return box.Render(msg)
}

func (a *App) renderErrors() string {
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("203")).Padding(0, 1)
	if len(a.lastErrors) == 0 {
		return ""
	}
	lines := append([]string{"Errors:"}, a.lastErrors...)
	return box.Render(strings.Join(lines, "\n"))
}

func (a *App) renderHelp() string {
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("245")).Padding(0, 1)
	text := "Keys: ↑/↓ or j/k move • a actions • e errors • f filter • s sort • o order • space select • enter details • r refresh • g/G home/end • ? help • q quit"
	return box.Render(text)
}

func (a *App) statusSummary() string {
	parts := []string{}
	parts = append(parts, fmt.Sprintf("filter=%s", a.statusFilter))
	parts = append(parts, fmt.Sprintf("sort=%s", a.sortBy))
	if a.sortAsc {
		parts = append(parts, "asc")
	} else {
		parts = append(parts, "desc")
	}
	if sel := a.selectedCount(); sel > 0 {
		parts = append(parts, fmt.Sprintf("selected=%d", sel))
	}
	if len(a.lastErrors) > 0 {
		parts = append(parts, fmt.Sprintf("errors=%d (e)", len(a.lastErrors)))
	}
	if a.statusLine != "" {
		parts = append(parts, "| "+a.statusLine)
	}
	return strings.Join(parts, "  ")
}

func (a *App) renderActions() string {
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("178")).Padding(0, 1)
	items := []string{"Delete", "Restore", "Move", "Mark Uploaded", "Cancel"}
	lines := []string{"Actions:"}
	for i, it := range items {
		cursor := "  "
		if a.actionIndex == i {
			cursor = "➜ "
		}
		lines = append(lines, cursor+it)
	}
	content := strings.Join(lines, "\n")
	return box.Render(content)
}
func (a *App) loadArchivesCmd() tea.Cmd {
	return func() tea.Msg {
		cfg, _ := config.Load()
		mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
		if err != nil {
			return archivesMsg{err: err}
		}
		defer mgr.Close()
		arcs, err := mgr.List()
		if err != nil {
			return archivesMsg{err: err}
		}
		return archivesMsg{archives: arcs}
	}
}

// Helpers: selection, sorting, filtering, rows
func (a *App) selectedCount() int {
	c := 0
	for _, v := range a.selected {
		if v {
			c++
		}
	}
	return c
}

func (a *App) cycleSortField() {
	switch a.sortBy {
	case "created":
		a.sortBy = "name"
	case "name":
		a.sortBy = "size"
	case "size":
		a.sortBy = "status"
	default:
		a.sortBy = "created"
	}
}

func (a *App) cycleStatusFilter() {
	switch a.statusFilter {
	case "all":
		a.statusFilter = "managed"
	case "managed":
		a.statusFilter = "external"
	case "external":
		a.statusFilter = "missing"
	case "missing":
		a.statusFilter = "deleted"
	default:
		a.statusFilter = "all"
	}
}

func (a *App) filtered() []*storage.Archive {
	list := a.archives
	switch a.statusFilter {
	case "managed":
		res := make([]*storage.Archive, 0, len(list))
		for _, arc := range list {
			if arc.Managed {
				res = append(res, arc)
			}
		}
		list = res
	case "external":
		res := make([]*storage.Archive, 0, len(list))
		for _, arc := range list {
			if !arc.Managed {
				res = append(res, arc)
			}
		}
		list = res
	case "missing":
		res := make([]*storage.Archive, 0, len(list))
		for _, arc := range list {
			if arc.Status == "missing" {
				res = append(res, arc)
			}
		}
		list = res
	case "deleted":
		res := make([]*storage.Archive, 0, len(list))
		for _, arc := range list {
			if arc.Status == "deleted" {
				res = append(res, arc)
			}
		}
		list = res
	}
	return list
}

func (a *App) refreshRows() {
	list := a.filtered()
	switch a.sortBy {
	case "name":
		sort.Slice(list, func(i, j int) bool {
			if a.sortAsc {
				return list[i].Name < list[j].Name
			}
			return list[i].Name > list[j].Name
		})
	case "size":
		sort.Slice(list, func(i, j int) bool {
			if a.sortAsc {
				return list[i].Size < list[j].Size
			}
			return list[i].Size > list[j].Size
		})
	case "status":
		sort.Slice(list, func(i, j int) bool {
			if a.sortAsc {
				return list[i].Status < list[j].Status
			}
			return list[i].Status > list[j].Status
		})
	default:
		sort.Slice(list, func(i, j int) bool {
			if a.sortAsc {
				return list[i].Created.Before(list[j].Created)
			}
			return list[i].Created.After(list[j].Created)
		})
	}
	rows := a.rowsFromArchives(list)
	prev := a.tbl.Cursor()
	a.tbl.SetRows(rows)
	if prev >= len(rows) {
		a.tbl.GotoTop()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
