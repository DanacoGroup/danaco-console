#!/usr/bin/env python3
"""Składa treść warunków licencji z dokumentu źródłowego do pliku treści.

Wynikiem jest `zasoby/tresci/licencja.js` — skrypt wpisujący fragment HTML do
rejestru treści `window.DanacoTresci`, skąd bierze go składnik `dokument`
w kroku 2 instalatora. Skrypt, a nie plik do pobrania, bo okno bywa otwierane
wprost z dysku — przeglądarka blokuje wtedy pobieranie plików towarzyszących.

Wymiana wydania licencji polega na podmianie dokumentu źródłowego i ponownym
uruchomieniu tego skryptu; nic w oknie nie jest pisane ręcznie.

    python3 zbuduj-licencje.py <sciezka-do-LICENSE.md>
"""
import re, subprocess, sys, pathlib

import pathlib as _p
ZRODLO = sys.argv[1] if len(sys.argv) > 1 else \
    str(_p.Path(__file__).parent / 'zasoby' / 'tresci' / 'licencja-2.1.md')
CEL = pathlib.Path(__file__).parent / 'zasoby' / 'tresci' / 'licencja.js'

tekst = pathlib.Path(ZRODLO).read_text(encoding='utf-8')

# Metryka i spis treści są zapisem wewnętrznym redakcji i nawigacją, nie treścią
# warunków. W oknie instalatora czytelnik zaczyna od rozdziału 1 — spis 100 pozycji
# na wejściu przesłaniałby to, co ma przeczytać.
poczatek = tekst.find('## 1. Postanowienia wstępne')
if poczatek < 0:
    sys.exit('Brak rozdziału 1 w dokumencie źródłowym — nie wiem, gdzie zaczyna się akt.')
tresc = tekst[poczatek:]

html = subprocess.run(
    ['pandoc', '-f', 'gfm+gfm_auto_identifiers', '-t', 'html', '--wrap=none'],
    input=tresc, capture_output=True, text=True, check=True).stdout

# Zestawienia bywają szersze niż kolumna tekstu — przewijają się we własnym polu.
html = html.replace('<table>', '<div class="dn-zestawienie-pole"><table>')
html = html.replace('</table>', '</table></div>')

# Odsyłacze do pozostałych dokumentów produktu nie mają dokąd prowadzić w oknie
# instalatora — zostaje sam tekst, bez odsyłacza donikąd.
html = re.sub(r'<a href="(?!#)[^"]*">([^<]*)</a>', r'\1', html)

# Treść wchodzi w literał szablonowy, więc znaki o własnym znaczeniu w takim
# literale muszą ustąpić — kolejność ucieczek jest wiążąca.
uciekniete = html.strip().replace('\\', '\\\\').replace('`', '\\`').replace('${', '\\${')

CEL.write_text(
    '/* ============================================================================\n'
    '   TREŚĆ UMOWY LICENCYJNEJ — wynik `zbuduj-licencje.py` ze źródła\n'
    '   `zasoby/tresci/licencja-2.1.md`. Fragment, nie dokument: wstawia go\n'
    '   składnik `dokument` do pola przewijanego w kroku 2.\n'
    '\n'
    '   Nie edytuj tu niczego — zmiana wchodzi w źródle i przechodzi przez\n'
    '   generator. Treść stoi w skrypcie, a nie w osobnym pliku do pobrania, bo\n'
    '   okno bywa otwierane wprost z dysku, a przeglądarka blokuje wtedy\n'
    '   pobieranie plików towarzyszących.\n'
    '   ============================================================================ */\n'
    'window.DanacoTresci = window.DanacoTresci || {};\n'
    'window.DanacoTresci.licencja = `\n' + uciekniete + '\n`;\n',
    encoding='utf-8')

cele = set(re.findall(r'id="([^"]+)"', html))
odsyl = set(re.findall(r'href="#([^"]+)"', html))
print('%s — wstawiono %d znaków, %d rozdziałów, %d odsyłaczy spisu, %d donikąd'
      % (CEL.name, len(html), html.count('<h2'), len(odsyl), len(odsyl - cele)))
