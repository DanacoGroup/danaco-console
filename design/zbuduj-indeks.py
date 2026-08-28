#!/usr/bin/env python3
# Skrypt buduje INDEKS.html ze skanu bieżącego stanu drzewa katalogu design na
# potrzeby przeglądu prototypów i opracowań.
"""Buduje INDEKS.html ze skanu drzewa design/. Katalog powstaje z rzeczywistości,
   więc nie może się z nią rozminąć."""
import os, re, html, sys, subprocess
D = '/home/ubuntu/robocze/prototypy/design'
S = os.path.join(D, 'zasoby', 'indeks')   # szablony katalogu
os.chdir(D)

def tytul(p):
    try: t = open(p, encoding='utf-8', errors='ignore').read(4000)
    except: return None
    m = re.search(r'<title>(.*?)</title>', t, re.S)
    if m: return html.unescape(re.sub(r'\s+', ' ', m.group(1))).strip()
    m = re.search(r'^#\s+(.+)$', t, re.M)
    return m.group(1).strip() if m else None

def miara(p):
    w = sum(1 for _ in open(p, encoding='utf-8', errors='ignore'))
    kb = os.path.getsize(p) / 1024
    return f'{w} w. · {kb:.0f} kB'

def klikalny(p):
    if not p.endswith('.html'): return False
    t = open(p, encoding='utf-8', errors='ignore').read()
    return bool(re.search(r'data-(idz|przelacz|krok-dalej|ekran)', t))

def karta(sciezka, plakietki=()):
    t = tytul(sciezka) or os.path.basename(sciezka)
    nazwa = os.path.basename(sciezka)
    typ = 'HTML' if sciezka.endswith('.html') else ('MD' if sciezka.endswith('.md') else 'CSS')
    pl = [f'<span class="dn-plakietka dn-plakietka--informacja">{typ}</span>' if typ=='HTML'
          else f'<span class="dn-plakietka">{typ}</span>']
    if klikalny(sciezka):
        pl.append('<span class="dn-plakietka dn-plakietka--sygnal"><span class="ix-kropka"></span>klikalny</span>')
    for p in plakietki:
        pl.append(f'<span class="dn-plakietka dn-plakietka--sukces">{html.escape(p)}</span>')
    szuk = html.escape((t + ' ' + nazwa).lower(), quote=True)
    return (f'<a class="dn-karta dn-karta--klikalna ix-karta" href="{html.escape(sciezka)}"'
            f' data-filtr-pozycja data-filtr-tekst="{szuk}">'
            f'<div class="ix-karta-gora">{"".join(pl)}</div>'
            f'<div class="ix-karta-tytul">{html.escape(t)}</div>'
            f'<div class="ix-karta-meta"><code>{html.escape(nazwa)}</code>'
            f'<span>{miara(sciezka)}</span></div></a>')

def zbierz(kat, rozsz=('.html','.md')):
    out = []
    for r, _, f in os.walk(kat):
        for x in sorted(f):
            if x.endswith(rozsz): out.append(os.path.join(r, x).replace('./',''))
    return sorted(out)

IK = ('<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" '
      'stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">'
      '<path d="M4 6h16M4 12h16M4 18h10"/></svg>')

def sekcja(idn, tytul_s, opis, pozycje, wyroz=False):
    if not pozycje: return '', None
    karty = ''.join(pozycje)
    kl = ' ix-sekcja--wyrozniona' if wyroz else ''
    return (f'<section class="ix-sekcja{kl}" id="{idn}">'
            f'<header class="ix-sekcja-naglowek"><span class="ix-sekcja-ikona">{IK}</span>'
            f'<div><h2>{html.escape(tytul_s)}</h2><p>{html.escape(opis)}</p></div>'
            f'<span class="ix-sekcja-licznik">{len(pozycje)}</span></header>'
            f'<div class="ix-siatka">{karty}</div></section>'), (idn, tytul_s, len(pozycje))

sekcje, spis = [], []

# Sekcja przedmiotu bieżącej pracy pokazuje prototypy oznaczone plakietką
# postępu, wyróżnione wizualnie na czele katalogu.
prace = [karta('05-okna/platformowe/instalator.html', ('w pracy',)),
         karta('05-okna/przeplyw/przeplyw-wejscia.html', ('w pracy',))]
s1, w1 = sekcja('w-pracy', 'Przedmiot bieżącej pracy',
    'Dwa prototypy przebudowane w sierpniu 2026: instalacja aplikacji oraz uruchomienie, rejestracja i logowanie.',
    prace, wyroz=True)
sekcje.append(s1); spis.append(w1)

# Sekcja pozostałych okien zbiera pliki katalogu okien pominięte w sekcji
# bieżącej pracy, posortowane alfabetycznie.
poz = [p for p in zbierz('05-okna') if p.endswith('.html')
       and 'instalator.html' not in p and 'przeplyw-wejscia.html' not in p]
s2, w2 = sekcja('okna', 'Prototypy okien platformy',
    'Okna platformowe, przepływy, środowiska i moduły — po jednym pliku na okno.',
    [karta(p) for p in poz])
sekcje.append(s2); spis.append(w2)

# Sekcja norm i wykazów zbiera opracowania tekstowe katalogu okien,
# obowiązujące przy budowie kolejnych prototypów.
normy = [p for p in zbierz('05-okna') if p.endswith('.md')]
s3, w3 = sekcja('normy', 'Normy i wykazy',
    'Zasady obowiązujące przy budowie okien oraz rejestr długu składników.',
    [karta(p) for p in normy])
sekcje.append(s3); spis.append(w3)

for kat, idn, tyt, op in [
    ('01-dokumentacja-md','dokumentacja','Opracowania merytoryczno-techniczne','System projektowy opisany tekstem.'),
    ('02-dokumentacja-html','dokumentacja-html','Opracowania merytoryczno-graficzne','Te same opracowania w postaci przeglądarkowej.'),
    ('03-marka','marka','Marka, księga znaku i portfolio logotypu','Konstrukcja znaku, zastosowania, typografia i głos.'),
    ('04-portfolio','portfolio','Portfolio design identity','Plansze portfolio osadzone w rzeczywistych danych platformy.'),
    ('_prace','prace','Próby i zestawienia robocze','Pliki pomocnicze — nie wchodzą do produktu.'),
]:
    if not os.path.isdir(kat): continue
    s, w = sekcja(idn, tyt, op, [karta(p) for p in zbierz(kat)])
    sekcje.append(s); spis.append(w)

# Blok metryk liczy opracowania, linie i okna drzewa design oraz składniki
# biblioteki wyprowadzone z arkuszy stylu.
wszystkie = [p for w in [zbierz('05-okna'), zbierz('01-dokumentacja-md'), zbierz('02-dokumentacja-html'),
                          zbierz('03-marka'), zbierz('04-portfolio'), zbierz('_prace')] for p in w]
linii = sum(sum(1 for _ in open(p, encoding='utf-8', errors='ignore')) for p in wszystkie)
ARK = "css/fundament.css css/komponenty.css menu.css rama.css izolacja.css karty-okna.css okna-modalne.css okno-robocze.css panel-sesji.css pasek-okna.css przedsionek.css stanowisko.css stany.css".split()
skl = subprocess.run(['bash','-c',
    "grep -ohE '^\\.dn-[a-z0-9-]+' " + " ".join('zasoby/'+a for a in ARK) + " | sort -u | wc -l"],
    capture_output=True, text=True).stdout.strip()

metryki = ''.join(f'<div class="ix-metryka"><b>{b}</b><span>{s}</span></div>' for b, s in [
    (len(wszystkie), 'opracowań'), (f'{linii:,}'.replace(',', ' '), 'linii'),
    (len([p for p in wszystkie if p.startswith('05-okna') and p.endswith('.html')]), 'okien'),
    (skl, 'składników biblioteki')])

boczna = ''.join(f'<a href="#{i}">{html.escape(t)}<span>{n}</span></a>' for i, t, n in spis if i)

styl = open(f'{S}/ix-style.html', encoding='utf-8').read()
pasek = open(f'{S}/ix-pasek.html', encoding='utf-8').read()
skrypty = open(f'{S}/ix-skrypty.html', encoding='utf-8').read()

doc = f'''<!doctype html>
<html lang="pl" data-theme="dark">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Danaco Console — katalog warstwy projektowej</title>
<link rel="icon" href="zasoby/marka/favicon/favicon.svg" type="image/svg+xml">
<link rel="stylesheet" href="zasoby/zetony/fonty.css">
<link rel="stylesheet" href="zasoby/zetony/zetony.css">
<link rel="stylesheet" href="zasoby/css/fundament.css">
<link rel="stylesheet" href="zasoby/css/komponenty.css">
<link rel="stylesheet" href="zasoby/prototyp.css">
{styl}
<style>
  .ix-sekcja--wyrozniona {{ padding: var(--dn-od-5); border-radius: var(--dn-r-lg);
    background: var(--dn-sygnal-tlo); border: 1px solid var(--dn-sygnal-obrys); }}
  .ix-szukaj {{ margin-bottom: var(--dn-od-6); }}
</style>
</head>
<body>
{pasek}
<div class="ix-uklad">
  <nav class="ix-boczna" data-spis aria-label="Spis sekcji">
    <div class="pt-etykieta" style="padding:0 var(--dn-od-3) var(--dn-od-2)">Sekcje</div>
    {boczna}
  </nav>
  <main>
    <header class="ix-naglowek">
      <svg class="ix-godlo" viewBox="0 0 96 96" role="img" aria-label="Danaco Console">
        <path fill="currentColor" d="M12 26 H24 L44 48 L24 70 H12 L32 48 Z"/>
        <path fill="currentColor" d="M40 26 H52 L72 48 L52 70 H40 L60 48 Z"/>
        <circle cx="83" cy="63.5" r="6.5" fill="var(--dn-sygnal-500)"/></svg>
      <div>
        <h1><b>Warstwa projektowa</b>Danaco Console</h1>
        <p>Katalog całości opracowań projektowych. Powstaje ze skanu drzewa, więc odpowiada
        stanowi plików, a nie zapisowi sprzed zmian.</p>
      </div>
    </header>
    <div class="ix-metryki">{metryki}</div>
    <div class="ix-szukaj">
      <input class="dn-pole-kontrolka" id="ix-filtr" type="search" placeholder="Szukaj w katalogu…"
             aria-label="Szukaj w katalogu" style="width:100%">
    </div>
    {''.join(sekcje)}
  </main>
</div>
{skrypty}
</body>
</html>
'''
open('INDEKS.html','w',encoding='utf-8').write(doc)
print('INDEKS.html zbudowany:', len(doc), 'znaków |', sum(n for _,_,n in spis if _), 'pozycji w', len(spis), 'sekcjach')
