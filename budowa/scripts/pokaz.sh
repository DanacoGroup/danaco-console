#!/usr/bin/env bash
# Podgląd stanu Danaco Console: skrypt składa stan drzewa, sprawdza, czy klient się buduje, i podnosi serwer klienta oraz rdzeń, aby stan dało się obejrzeć w przeglądarce na żądanie o każdej porze.

set -u

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BUDOWA="$KORZEN/budowa"
PORT_KLIENTA=5173
PORT_RDZENIA=17870

if [ "${1:-}" = "--stop" ]; then
  echo "Zatrzymuję serwery podglądu…"
  taskkill //FI "WINDOWTITLE eq *vite*" //F >/dev/null 2>&1 || true
  for p in "$PORT_KLIENTA" "$PORT_RDZENIA"; do
    pid=$(netstat -ano 2>/dev/null | grep -E "LISTENING" | grep ":$p " | awk '{print $NF}' | head -1)
    [ -n "${pid:-}" ] && taskkill //PID "$pid" //F >/dev/null 2>&1 && echo "  zatrzymany port $p"
  done
  exit 0
fi

echo "══ Stan aktualny"
printf '  zapisów w repozytorium: %s\n' "$(git -C "$KORZEN" rev-list --count HEAD 2>/dev/null)"
printf '  plików w budowa:        %s\n' "$(find "$BUDOWA" -type f | grep -v node_modules | grep -v '/target/' | wc -l)"

echo
echo "══ Co już istnieje"
for skladnik in \
  "server/internal/transport:serwer WebSocket" \
  "server/internal/core:rdzeń — dispatch" \
  "server/internal/session:sesje i okna komunikacji" \
  "server/internal/models:rejestr kanałów modelu" \
  "server/internal/injection:kanał Claude Code" \
  "client/src/motyw:żetony systemu wizualnego" \
  "client/src/komponenty:biblioteka komponentów" \
  "client/src/sterowanie:sterowanie per okno" \
  "desktop:powłoka natywna Tauri"
do
  sciezka="${skladnik%%:*}"; opis="${skladnik#*:}"
  n=$(find "$BUDOWA/$sciezka" -type f 2>/dev/null | grep -v '/target/' | wc -l)
  if [ "$n" -gt 0 ]; then printf '  \033[32m✓\033[0m %-34s %s plików\n' "$opis" "$n"
  else                    printf '  \033[90m·\033[0m %-34s —\n' "$opis"; fi
done

echo
echo "══ Czy klient się buduje"
if (cd "$BUDOWA/klient" && npx tsc --noEmit >/dev/null 2>&1); then
  echo "  ✓ kontrola typów przechodzi"
else
  echo "  ✗ kontrola typów NIE przechodzi — okno może nie odzwierciedlać kodu"
fi

if (cd "$BUDOWA/klient" && npm run budowanie >/tmp/pokaz-build.txt 2>&1); then
  rozmiar=$(grep -oE 'index-[A-Za-z0-9_-]+\.js *[0-9.]+ kB' /tmp/pokaz-build.txt | head -1 || echo '?')
  echo "  ✓ pakiet zbudowany: $rozmiar"
else
  echo "  ✗ budowanie klienta nie przechodzi"; tail -6 /tmp/pokaz-build.txt | sed 's/^/      /'
fi

echo
echo "══ Uruchamiam podgląd"

# Serwer rozwojowy klienta.
# Uruchamiamy przez `start`, a nie przez `&` — proces w tle podpięty do tej samej
# konsoli zatrzymuje skrypt do czasu własnego zakończenia, a serwer nie kończy się
# nigdy.
if netstat -ano 2>/dev/null | grep -q ":$PORT_KLIENTA .*LISTENING"; then
  echo "  serwer klienta już działa na $PORT_KLIENTA"
else
  ( cd "$BUDOWA/klient" && start //B cmd //c "npm run dev > %TEMP%\\pokaz-vite.txt 2>&1" ) >/dev/null 2>&1
  for _ in 1 2 3 4 5 6 7 8; do
    netstat -ano 2>/dev/null | grep -q ":$PORT_KLIENTA .*LISTENING" && break
    sleep 1
  done
  netstat -ano 2>/dev/null | grep -q ":$PORT_KLIENTA .*LISTENING" \
    && echo "  serwer klienta uruchomiony na $PORT_KLIENTA" \
    || echo "  serwer klienta nie wstał — zajrzyj do %TEMP%\\pokaz-vite.txt"
fi

# Rdzeń — wyłącznie gdy istnieje transport i gdy da się zbudować.
# Budujemy do pliku wykonywalnego zamiast `go run`: `go run` trzyma proces
# nadrzędny, przez co skrypt nie oddaje sterowania.
if [ -d "$BUDOWA/server/internal/transport" ]; then
  if netstat -ano 2>/dev/null | grep -q ":$PORT_RDZENIA .*LISTENING"; then
    echo "  rdzeń już działa na $PORT_RDZENIA"
  elif (cd "$BUDOWA" && go build -o /tmp/danaco-console.exe ./server/cmd/danaco-console 2>/dev/null); then
    ( start //B cmd //c "/tmp/danaco-console.exe > %TEMP%\\pokaz-rdzen.txt 2>&1" ) >/dev/null 2>&1
    for _ in 1 2 3 4 5; do
      netstat -ano 2>/dev/null | grep -q ":$PORT_RDZENIA .*LISTENING" && break
      sleep 1
    done
    netstat -ano 2>/dev/null | grep -q ":$PORT_RDZENIA .*LISTENING" \
      && echo "  rdzeń uruchomiony na $PORT_RDZENIA" \
      || echo "  rdzeń nie nasłuchuje — transport może jeszcze nie być podpięty"
  else
    echo "  rdzeń: nie buduje się — okno działa bez połączenia"
  fi
else
  echo "  rdzeń: brak transportu — okno działa bez połączenia"
fi

echo
echo "══ Otwieram"
start "http://localhost:$PORT_KLIENTA/" 2>/dev/null || true
echo "  okno aplikacji:      http://localhost:$PORT_KLIENTA/"

PODGLAD="$KORZEN/opracowania/system-wizualny/podglad.html"
if [ -f "$PODGLAD" ]; then
  start "$PODGLAD" 2>/dev/null || true
  echo "  system wizualny:     opracowania/system-wizualny/podglad.html"
fi

echo
echo "  zatrzymanie: bash budowa/scripts/pokaz.sh --stop"
