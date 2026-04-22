package integration

import (
	"strings"
	"testing"
)

func TestLocaleHelpRootUkrainian(t *testing.T) {
	bin := buildBinary(t, "")

	stdout, stderr, code := runBinaryWithEnv(t, bin, []string{"LANG=uk_UA.UTF-8"}, "help")

	if code != 0 {
		t.Errorf("exit code = %d, want 0 (stderr=%q)", code, stderr)
	}

	for _, want := range []string{
		"mytets — легкий інструмент командного рядка",
		"Використання:",
		"[команда]",
		"Доступні команди:",
		"Показати список випадкових фраз",
		"Показати одну випадкову фразу",
		"Вивести версію програми та завершити",
		"Прапори:",
		"довідка для mytets",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("Ukrainian root help missing %q\nGot:\n%s", want, stdout)
		}
	}
}

func TestLocaleHelpOneUkrainian(t *testing.T) {
	bin := buildBinary(t, "")

	stdout, stderr, code := runBinaryWithEnv(t, bin, []string{"LANG=uk_UA.UTF-8"}, "one", "--help")

	if code != 0 {
		t.Errorf("exit code = %d, want 0 (stderr=%q)", code, stderr)
	}

	for _, want := range []string{
		"Команда one виводить випадкову фразу",
		"Використання:",
		"[прапори]",
		"Прапори:",
		"довідка для one",
		"Формат виводу:",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("Ukrainian one --help missing %q\nGot:\n%s", want, stdout)
		}
	}
}

func TestLocaleHelpListUkrainian(t *testing.T) {
	bin := buildBinary(t, "")

	stdout, stderr, code := runBinaryWithEnv(t, bin, []string{"LANG=uk_UA.UTF-8"}, "list", "--help")

	if code != 0 {
		t.Errorf("exit code = %d, want 0 (stderr=%q)", code, stderr)
	}

	for _, want := range []string{
		"Команда list виводить кілька унікальних випадкових фраз",
		"Використання:",
		"[прапори]",
		"Прапори:",
		"довідка для list",
		"Кількість фраз для повернення",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("Ukrainian list --help missing %q\nGot:\n%s", want, stdout)
		}
	}
}

func TestLocaleHelpVersionUkrainian(t *testing.T) {
	bin := buildBinary(t, "")

	stdout, stderr, code := runBinaryWithEnv(t, bin, []string{"LANG=uk_UA.UTF-8"}, "version", "--help")

	if code != 0 {
		t.Errorf("exit code = %d, want 0 (stderr=%q)", code, stderr)
	}

	for _, want := range []string{
		"Вивести версію програми та завершити",
		"Використання:",
		"[прапори]",
		"Прапори:",
		"довідка для version",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("Ukrainian version --help missing %q\nGot:\n%s", want, stdout)
		}
	}
}

func TestLocaleHelpRootEnglish(t *testing.T) {
	bin := buildBinary(t, "")

	stdout, stderr, code := runBinaryWithEnv(t, bin, []string{"LANG=en_US.UTF-8"}, "help")

	if code != 0 {
		t.Errorf("exit code = %d, want 0 (stderr=%q)", code, stderr)
	}

	for _, want := range []string{
		"mytets — a lightweight CLI tool",
		"Usage:",
		"[command]",
		"Available Commands:",
		"Flags:",
		"help for mytets",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("English root help missing %q\nGot:\n%s", want, stdout)
		}
	}
}
