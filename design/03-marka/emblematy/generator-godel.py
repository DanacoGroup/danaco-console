#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# Skrypt wytwarza komplet godeł pochodnych marki: emblematy w postaci wektorowej
# i rastrowej, ikonę aplikacji oraz pliki favicon, z jednej zatwierdzonej geometrii.
"""
Danaco Console — generator systemu godeł pochodnych.

Wytwarza z ZATWIERDZONEJ geometrii (bez modyfikacji krzywych):
  * warianty rozmiarowe czterech emblematów środowisk (16/20/24/32/48/64 px, SVG),
  * warianty markowe emblematów (kropka w błękicie sygnałowym) na jasnym i ciemnym,
  * rastry PNG emblematów (24/48/96/192 px × jasny/ciemny),
  * komplet ikony aplikacji (PNG 32…1024 + wariant maskowalny + ICO wielorozmiarowy),
  * komplet faviconu (SVG adaptacyjny + kafel SVG + PNG 16/32/48 + ICO + apple-touch-icon),
  * site.webmanifest i snippet <head>.

Uruchomienie:  python3 generator-godel.py
Zależności:    cairosvg, Pillow
"""

import io
import os
import struct

import cairosvg
from PIL import Image

# Katalogi wyjściowe powstają przed wytworzeniem plików, osobno dla emblematów,
# ikony aplikacji oraz kompletu favicon.
BAZA = os.path.dirname(os.path.abspath(__file__))
KAT_SVG = os.path.join(BAZA, "svg")
KAT_SVG_MARKA = os.path.join(KAT_SVG, "warianty")
KAT_PNG = os.path.join(BAZA, "png")
KAT_IKONA = os.path.join(BAZA, "ikona-aplikacji")
KAT_FAVICON = os.path.join(BAZA, "favicon")
for k in (KAT_SVG, KAT_SVG_MARKA, KAT_PNG, KAT_IKONA, KAT_FAVICON):
    os.makedirs(k, exist_ok=True)

# Barwy odpowiadają żetonom systemu projektowego, lecz w plikach wektorowych
# zapisuje się je dosłownie, ponieważ plik wychodzi poza arkusz stylów.
INK_JASNY = "#16181D"   # --dn-tekst (motyw jasny)  = szary-900
INK_CIEMNY = "#ECECEC"  # --dn-tekst (motyw ciemny) = szary-100
DOT_JASNY = "#3B6FE0"   # --dn-kropka (motyw jasny)  = sygnal-500
DOT_CIEMNY = "#5C8CEC"  # --dn-kropka (motyw ciemny) = sygnal-400
SYGNAL_JASNY = "#2457C9"   # --dn-sygnal (motyw jasny)  = sygnal-600
SYGNAL_CIEMNY = "#8FB2F5"  # --dn-sygnal (motyw ciemny) = sygnal-300
GRUNT = "#111317"       # --dn-rama = szary-925 — grunt kafla i ikony aplikacji

# Zatwierdzona geometria sygnetu opisana jest na siatce dziewięćdziesięciu sześciu
# jednostek i stanowi jedyne źródło kształtu dla wszystkich plików pochodnych.
GROT_1 = "M12 26 H24 L44 48.0 L24 70 H12 L32 48.0 Z"
GROT_2 = "M40 26 H52 L72 48.0 L52 70 H40 L60 48.0 Z"
KROPKA = ("83", "63.5", "6.5")          # cx, cy, r
GROT_UPR = "M18 18 H35 L64 48.0 L35 78 H18 L45 48.0 Z"
KROPKA_UPR = ("79", "69", "9")

# ── ZATWIERDZONA geometria czterech emblematów (siatka 24×24, obrys 1,75) ───
# Klucz: (elementy-obrysowe, kropka(cx, cy, r))
EMBLEMATY = {
    "srodowisko-talkin": {
        "nazwa": "TalkIn",
        "obrys": [
            '<path d="M20.5 5.5a2.5 2.5 0 0 0-2.5-2.5H6A2.5 2.5 0 0 0 3.5 5.5v8A2.5 2.5 0 0 0 6 16h1v4l4.4-4H18a2.5 2.5 0 0 0 2.5-2.5Z"/>',
            '<path d="M7.5 7.5h9"/>',
            '<path d="M7.5 10.8h5"/>',
        ],
        "kropka": ("16.2", "10.8", "1.5"),
        # Wariant szesnastopikselowy: górna linia tekstu usunięta, kropka powiększona.
        "obrys16": [
            '<path d="M20.5 5.5a2.5 2.5 0 0 0-2.5-2.5H6A2.5 2.5 0 0 0 3.5 5.5v8A2.5 2.5 0 0 0 6 16h1v4l4.4-4H18a2.5 2.5 0 0 0 2.5-2.5Z"/>',
            '<path d="M7.5 9.4h4.6"/>',
        ],
        "kropka16": ("16.2", "9.4", "1.9"),
        "obrys16w": "1.9",
    },
    "srodowisko-workspace": {
        "nazwa": "WorkSpace",
        "obrys": [
            '<rect x="3.5" y="3.5" width="7.2" height="7.2" rx="2"/>',
            '<rect x="13.3" y="3.5" width="7.2" height="7.2" rx="2"/>',
            '<rect x="3.5" y="13.3" width="7.2" height="7.2" rx="2"/>',
            '<rect x="13.3" y="13.3" width="7.2" height="7.2" rx="2"/>',
        ],
        "kropka": ("16.9", "16.9", "1.6"),
        # Wariant szesnastopikselowy: obrys czwartego modułu ustępuje samej kropce.
        "obrys16": [
            '<rect x="3.55" y="3.55" width="6.9" height="6.9" rx="1.4"/>',
            '<rect x="13.55" y="3.55" width="6.9" height="6.9" rx="1.4"/>',
            '<rect x="3.55" y="13.55" width="6.9" height="6.9" rx="1.4"/>',
        ],
        "kropka16": ("17", "17", "3"),
        "obrys16w": "1.9",
    },
    "srodowisko-codestudio": {
        "nazwa": "CodeStudio",
        "obrys": [
            '<rect x="3" y="4" width="18" height="16" rx="2.5"/>',
            '<path d="M7 9.2l3.1 2.8L7 14.8"/>',
            '<path d="M12.6 15.4h4"/>',
        ],
        "kropka": ("18.4", "9.9", "1.5"),
        # Wariant szesnastopikselowy: linia wyniku usunięta, ramka zwężona.
        "obrys16": [
            '<rect x="3.2" y="4.2" width="17.6" height="15.6" rx="2.2"/>',
            '<path d="M7.2 9.4l3.1 2.8-3.1 2.8"/>',
        ],
        "kropka16": ("17.4", "10", "1.85"),
        "obrys16w": "1.9",
    },
    "srodowisko-multitaskingai": {
        "nazwa": "MultitaskingAI",
        "obrys": [
            '<circle cx="12" cy="4.6" r="1.9"/>',
            '<circle cx="4.9" cy="16.4" r="1.9"/>',
            '<circle cx="19.1" cy="16.4" r="1.9"/>',
            '<path d="M12 9.4V6.5"/>',
            '<path d="M9.8 13.4l-3.3 1.9"/>',
            '<path d="M14.2 13.4l3.3 1.9"/>',
        ],
        "kropka": ("12", "12", "2.6"),
        # Wariant szesnastopikselowy: żaden element nie znika, węzły powiększone.
        "obrys16": [
            '<circle cx="12" cy="4.6" r="2"/>',
            '<circle cx="4.9" cy="16.4" r="2"/>',
            '<circle cx="19.1" cy="16.4" r="2"/>',
            '<path d="M12 9.1V6.6"/>',
            '<path d="M9.55 13.55l-2.95 1.8"/>',
            '<path d="M14.45 13.55l2.95 1.8"/>',
        ],
        "kropka16": ("12", "12", "2.9"),
        "obrys16w": "1.9",
    },
}

KOLEJNOSC = [
    "srodowisko-talkin",
    "srodowisko-workspace",
    "srodowisko-codestudio",
    "srodowisko-multitaskingai",
]

ROZMIARY_SVG = [16, 20, 24, 32, 48, 64]
# Grubość obrysu podana w jednostkach siatki dwudziestoczterojednostkowej rośnie
# poniżej dwudziestu czterech pikseli, kompensując optyczne zanikanie kreski.
OBRYS_DLA = {16: "1.9", 20: "1.9", 24: "1.75", 32: "1.75", 48: "1.75", 64: "1.75"}


def emblemat_svg(klucz, rozmiar, ink=None, dot=None):
    """Buduje SVG emblematu w zadanym rozmiarze.
    ink/dot = None → wariant interfejsowy (currentColor).
    ink/dot podane → wariant markowy (barwy dosłowne)."""
    e = EMBLEMATY[klucz]
    uproszczony = rozmiar <= 16
    obrys = e["obrys16"] if uproszczony else e["obrys"]
    kropka = e["kropka16"] if uproszczony else e["kropka"]
    szer = e["obrys16w"] if uproszczony else OBRYS_DLA[rozmiar]
    barwa_obrys = ink or "currentColor"
    barwa_kropka = dot or "currentColor"
    tytul = e["nazwa"]
    czesci = "".join(obrys)
    return (
        '<svg xmlns="http://www.w3.org/2000/svg" width="{r}" height="{r}" viewBox="0 0 24 24" '
        'role="img" aria-label="Środowisko {t}" fill="none" stroke="{so}" stroke-width="{w}" '
        'stroke-linecap="round" stroke-linejoin="round">'
        "{cz}"
        '<circle cx="{cx}" cy="{cy}" r="{cr}" fill="{sk}" stroke="none"/>'
        "</svg>\n"
    ).format(
        r=rozmiar, t=tytul, so=barwa_obrys, w=szer, cz=czesci,
        cx=kropka[0], cy=kropka[1], cr=kropka[2], sk=barwa_kropka,
    )


def zapisz(sciezka, tresc):
    with open(sciezka, "w", encoding="utf-8") as f:
        f.write(tresc)
    return sciezka


def png_z_svg(svg_tekst, sciezka, px, tlo=None):
    cairosvg.svg2png(
        bytestring=svg_tekst.encode("utf-8"),
        write_to=sciezka,
        output_width=px,
        output_height=px,
        background_color=tlo,
    )
    return sciezka


def png_bajty(svg_tekst, px):
    return cairosvg.svg2png(
        bytestring=svg_tekst.encode("utf-8"), output_width=px, output_height=px
    )


def zbuduj_ico(sciezka, pary):
    """pary = [(px, bajty_png)] — ICO z osobną grafiką dla każdego rozmiaru."""
    n = len(pary)
    naglowek = struct.pack("<HHH", 0, 1, n)
    wpisy, dane, offset = b"", b"", 6 + 16 * n
    for px, blob in pary:
        w = 0 if px >= 256 else px
        wpisy += struct.pack("<BBBBHHII", w, w, 0, 0, 1, 32, len(blob), offset)
        offset += len(blob)
        dane += blob
    with open(sciezka, "wb") as f:
        f.write(naglowek + wpisy + dane)
    return sciezka


# Emblematy w wariancie interfejsowym powstają w postaci wektorowej dla każdego
# rozmiaru, dziedzicząc barwę po otoczeniu zamiast nieść ją zapisaną w pliku.
licznik = {"svg": 0, "png": 0, "ico": 0, "inne": 0}

for klucz in KOLEJNOSC:
    for r in ROZMIARY_SVG:
        zapisz(os.path.join(KAT_SVG, "%s-%d.svg" % (klucz, r)), emblemat_svg(klucz, r))
        licznik["svg"] += 1

# Emblematy w wariantach barwnych niosą barwę wypaloną w pliku, dla rastrów oraz
# osadzeń niezdolnych dziedziczyć barwy po otoczeniu; cały emblemat ma jedną barwę.
WARIANTY_BARWNE = [
    ("jasny", INK_JASNY),           # na tle jasnym  — --dn-tekst (motyw jasny)
    ("ciemny", INK_CIEMNY),         # na tle ciemnym — --dn-tekst (motyw ciemny)
    ("sygnal-jasny", SYGNAL_JASNY),   # środowisko aktywne, motyw jasny  — --dn-sygnal
    ("sygnal-ciemny", SYGNAL_CIEMNY),  # środowisko aktywne, motyw ciemny — --dn-sygnal
]
for klucz in KOLEJNOSC:
    for przyrostek, barwa in WARIANTY_BARWNE:
        zapisz(
            os.path.join(KAT_SVG_MARKA, "%s-%s.svg" % (klucz, przyrostek)),
            emblemat_svg(klucz, 24, barwa, barwa),
        )
        licznik["svg"] += 1

# Emblematy w postaci rastrowej powstają jako pliki jednobarwne z tłem
# przezroczystym, w rozmiarach wymaganych przez osadzenia poza interfejsem.
for klucz in KOLEJNOSC:
    for px in (24, 48, 96, 192):
        for przyrostek, barwa in (("jasny", INK_JASNY), ("ciemny", INK_CIEMNY)):
            png_z_svg(
                emblemat_svg(klucz, 24, barwa, barwa),
                os.path.join(KAT_PNG, "%s-%d-%s.png" % (klucz, px, przyrostek)),
                px,
            )
            licznik["png"] += 1

# Ikona aplikacji powstaje z kompozycji zatwierdzonej, w której sygnet jest
# przeskalowany i odsunięty od krawędzi kafla o wartości podane niżej.
SYGNET_W_IKONIE = (
    '<g transform="translate(194.6,194.6) scale(6.6133)">'
    '<path fill="{ink}" d="{g1}"/><path fill="{ink}" d="{g2}"/>'
    '<circle cx="{cx}" cy="{cy}" r="{cr}" fill="{dot}"/></g>'
).format(ink=INK_CIEMNY, g1=GROT_1, g2=GROT_2,
         cx=KROPKA[0], cy=KROPKA[1], cr=KROPKA[2], dot=DOT_CIEMNY)

SYGNET_W_MASCE = (
    '<g transform="translate(245.8,245.8) scale(5.5467)">'
    '<path fill="{ink}" d="{g1}"/><path fill="{ink}" d="{g2}"/>'
    '<circle cx="{cx}" cy="{cy}" r="{cr}" fill="{dot}"/></g>'
).format(ink=INK_CIEMNY, g1=GROT_1, g2=GROT_2,
         cx=KROPKA[0], cy=KROPKA[1], cr=KROPKA[2], dot=DOT_CIEMNY)

# Wariant uproszczony sygnetu służy najmniejszym rastrom ikony aplikacji, gdzie
# pełna geometria zlewa się w plamę nieczytelną dla oka.
SYGNET_UPR_W_IKONIE = (
    '<g transform="translate(88,128) scale(8)">'
    '<path fill="{ink}" d="{g}"/>'
    '<circle cx="{cx}" cy="{cy}" r="{cr}" fill="{dot}"/></g>'
).format(ink=INK_CIEMNY, g=GROT_UPR,
         cx=KROPKA_UPR[0], cy=KROPKA_UPR[1], cr=KROPKA_UPR[2], dot=DOT_CIEMNY)

IKONA_SVG = (
    '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" role="img" '
    'aria-label="Danaco Console">\n'
    '  <rect width="1024" height="1024" rx="224" fill="%s"/>\n  %s\n</svg>\n'
    % (GRUNT, SYGNET_W_IKONIE)
)
IKONA_MASKA_SVG = (
    '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" role="img" '
    'aria-label="Danaco Console">\n'
    '  <rect width="1024" height="1024" rx="0" fill="%s"/>\n  %s\n</svg>\n'
    % (GRUNT, SYGNET_W_MASCE)
)
IKONA_KWADRAT_SVG = (
    '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" role="img" '
    'aria-label="Danaco Console">\n'
    '  <rect width="1024" height="1024" rx="0" fill="%s"/>\n  %s\n</svg>\n'
    % (GRUNT, SYGNET_W_IKONIE)
)
IKONA_UPR_SVG = (
    '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" role="img" '
    'aria-label="Danaco Console">\n'
    '  <rect width="1024" height="1024" rx="224" fill="%s"/>\n  %s\n</svg>\n'
    % (GRUNT, SYGNET_UPR_W_IKONIE)
)

zapisz(os.path.join(KAT_IKONA, "ikona-aplikacji.svg"), IKONA_SVG)
zapisz(os.path.join(KAT_IKONA, "ikona-maskowalna.svg"), IKONA_MASKA_SVG)
zapisz(os.path.join(KAT_IKONA, "ikona-kwadratowa.svg"), IKONA_KWADRAT_SVG)
zapisz(os.path.join(KAT_IKONA, "ikona-uproszczona.svg"), IKONA_UPR_SVG)
licznik["svg"] += 4

for px in (32, 64, 128, 180, 192, 256, 512, 1024):
    zrodlo = IKONA_UPR_SVG if px <= 32 else IKONA_SVG
    png_z_svg(zrodlo, os.path.join(KAT_IKONA, "ikona-%d.png" % px), px)
    licznik["png"] += 1

# Powłoka aplikacji oczekuje pliku o podwojonej gęstości pod nazwą wskazującą
# rozmiar podstawowy, czyli obrazu o boku dwustu pięćdziesięciu sześciu pikseli.
Image.open(os.path.join(KAT_IKONA, "ikona-256.png")).save(
    os.path.join(KAT_IKONA, "ikona-128@2x.png")
)
licznik["png"] += 1

for px in (192, 512):
    png_z_svg(IKONA_MASKA_SVG, os.path.join(KAT_IKONA, "ikona-maskowalna-%d.png" % px), px)
    licznik["png"] += 1

# Plik ikony aplikacji dla systemu Windows niesie rozmiary od szesnastu do
# dwustu pięćdziesięciu sześciu pikseli, z grafiką dobraną osobno do każdego.
pary = []
for px in (16, 24, 32, 48, 64, 128, 256):
    zrodlo = IKONA_UPR_SVG if px <= 32 else IKONA_SVG
    pary.append((px, png_bajty(zrodlo, px)))
zbuduj_ico(os.path.join(KAT_IKONA, "ikona-aplikacji.ico"), pary)
licznik["ico"] += 1

# Favicon wektorowy jest kopią pliku zatwierdzonego: tło przezroczyste, wariant
# uproszczony, barwy przełączane zapytaniem o motyw preferowany przez przeglądarkę.
FAVICON_SVG = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 96 96" role="img" aria-label="Danaco Console">
  <style>
    .znak { fill: %s; }
    .kropka { fill: %s; }
    @media (prefers-color-scheme: dark) {
      .znak { fill: %s; }
      .kropka { fill: %s; }
    }
  </style>
  <path class="znak" d="%s"/>
  <circle class="kropka" cx="%s" cy="%s" r="%s"/>
</svg>
""" % (INK_JASNY, DOT_JASNY, INK_CIEMNY, DOT_CIEMNY,
       GROT_UPR, KROPKA_UPR[0], KROPKA_UPR[1], KROPKA_UPR[2])
zapisz(os.path.join(KAT_FAVICON, "favicon.svg"), FAVICON_SVG)
licznik["svg"] += 1

# Kafel stoi na gruncie kryjącym, ponieważ nośnik rastrowy musi pozostać czytelny
# na pasku karty w motywie jasnym oraz ciemnym.
KAFEL_PELNY = """<svg xmlns="http://www.w3.org/2000/svg" width="96" height="96" viewBox="0 0 96 96" role="img" aria-label="Danaco Console">
<rect width="96" height="96" rx="21" fill="%s"/>
<g transform="translate(16.535,18.240) scale(0.62)"><path fill="%s" d="%s"/><path fill="%s" d="%s"/><circle cx="%s" cy="%s" r="%s" fill="%s"/></g>
</svg>
""" % (GRUNT, INK_CIEMNY, GROT_1, INK_CIEMNY, GROT_2,
       KROPKA[0], KROPKA[1], KROPKA[2], DOT_CIEMNY)

KAFEL_UPROSZCZONY = """<svg xmlns="http://www.w3.org/2000/svg" width="96" height="96" viewBox="0 0 96 96" role="img" aria-label="Danaco Console">
<rect width="96" height="96" rx="21" fill="%s"/>
<g transform="translate(9.84,13.44) scale(0.72)"><path fill="%s" d="%s"/><circle cx="%s" cy="%s" r="%s" fill="%s"/></g>
</svg>
""" % (GRUNT, INK_CIEMNY, GROT_UPR,
       KROPKA_UPR[0], KROPKA_UPR[1], KROPKA_UPR[2], DOT_CIEMNY)

zapisz(os.path.join(KAT_FAVICON, "favicon-kafel.svg"), KAFEL_PELNY)
zapisz(os.path.join(KAT_FAVICON, "favicon-kafel-uproszczony.svg"), KAFEL_UPROSZCZONY)
licznik["svg"] += 2

for px in (16, 32, 48):
    zrodlo = KAFEL_UPROSZCZONY if px <= 16 else KAFEL_PELNY
    png_z_svg(zrodlo, os.path.join(KAT_FAVICON, "favicon-%d.png" % px), px)
    licznik["png"] += 1

zbuduj_ico(
    os.path.join(KAT_FAVICON, "favicon.ico"),
    [
        (16, png_bajty(KAFEL_UPROSZCZONY, 16)),
        (32, png_bajty(KAFEL_PELNY, 32)),
        (48, png_bajty(KAFEL_PELNY, 48)),
    ],
)
licznik["ico"] += 1

# Ikona dotykowa systemu mobilnego jest kwadratem bez zaokrąglenia i bez kanału
# przezroczystości, ponieważ maskę nakłada sam system operacyjny.
att = os.path.join(KAT_FAVICON, "apple-touch-icon.png")
png_z_svg(IKONA_KWADRAT_SVG, att, 180)
Image.open(att).convert("RGB").save(att)
licznik["png"] += 1

# Ikony aplikacji instalowanej ze strony powstają z tego samego kompletu, co ikona
# aplikacji, w dwóch rozmiarach wymaganych przez manifest.
for px in (192, 512):
    png_z_svg(IKONA_SVG, os.path.join(KAT_FAVICON, "icon-%d.png" % px), px)
    licznik["png"] += 1
png_z_svg(IKONA_MASKA_SVG, os.path.join(KAT_FAVICON, "ikona-maskowalna-512.png"), 512)
licznik["png"] += 1

# Manifest aplikacji instalowanej ze strony oraz fragment nagłówka dokumentu
# powstają na końcu, gdy komplet plików graficznych stoi już na dysku.
MANIFEST = """{
  "id": "/",
  "name": "Danaco Console — AI Operating Environment",
  "short_name": "Danaco Console",
  "description": "Platforma operacyjna dla Operatora zarządzającego cyfrową organizacją: cztery środowiska, piętnaście modułów, wspólne okno komunikacji.",
  "lang": "pl",
  "dir": "ltr",
  "start_url": "/",
  "scope": "/",
  "display": "standalone",
  "display_override": ["window-controls-overlay", "standalone"],
  "orientation": "any",
  "background_color": "#0F0F0F",
  "theme_color": "#111317",
  "categories": ["productivity", "developer", "business"],
  "icons": [
    { "src": "/icon-192.png", "sizes": "192x192", "type": "image/png", "purpose": "any" },
    { "src": "/icon-512.png", "sizes": "512x512", "type": "image/png", "purpose": "any" },
    { "src": "/ikona-maskowalna-512.png", "sizes": "512x512", "type": "image/png", "purpose": "maskable" },
    { "src": "/favicon.svg", "sizes": "any", "type": "image/svg+xml", "purpose": "any" }
  ],
  "shortcuts": [
    { "name": "TalkIn", "url": "/talkin", "icons": [{ "src": "/emblematy/srodowisko-talkin-192.png", "sizes": "192x192", "type": "image/png" }] },
    { "name": "WorkSpace", "url": "/workspace", "icons": [{ "src": "/emblematy/srodowisko-workspace-192.png", "sizes": "192x192", "type": "image/png" }] },
    { "name": "CodeStudio", "url": "/codestudio", "icons": [{ "src": "/emblematy/srodowisko-codestudio-192.png", "sizes": "192x192", "type": "image/png" }] },
    { "name": "MultitaskingAI", "url": "/multitaskingai", "icons": [{ "src": "/emblematy/srodowisko-multitaskingai-192.png", "sizes": "192x192", "type": "image/png" }] }
  ]
}
"""
zapisz(os.path.join(KAT_FAVICON, "site.webmanifest"), MANIFEST)
licznik["inne"] += 1

SNIPPET = """<!-- ════════════════════════════════════════════════════════════════════════
     Danaco Console — komplet wpięć godła dokumentu (wklej do <head>)
     Kolejność jest znacząca: przeglądarka wspierająca SVG bierze pierwszy
     wpis i ignoruje ICO; starsza pobiera /favicon.ico ze ścieżki domyślnej.
     ════════════════════════════════════════════════════════════════════════ -->

<!-- 1. Nośnik podstawowy: SVG adaptacyjny (przełącza barwy z prefers-color-scheme) -->
<link rel="icon" href="/favicon.svg" type="image/svg+xml">

<!-- 2. Zapas rastrowy dla starszych przeglądarek: ICO 16/32/48 na kryjącym gruncie -->
<link rel="icon" href="/favicon.ico" sizes="16x16 32x32 48x48">

<!-- 3. iOS / iPadOS — kwadrat bez zaokrąglenia i bez kanału alfa (maskę nakłada system) -->
<link rel="apple-touch-icon" href="/apple-touch-icon.png">

<!-- 4. Manifest aplikacji internetowej (nazwa, barwy, ikony, skróty do środowisk) -->
<link rel="manifest" href="/site.webmanifest">

<!-- 5. Barwa paska systemowego — rama kokpitu w motywie ciemnym,
        papier roboczy w motywie jasnym -->
<meta name="theme-color" content="#F3F4F6" media="(prefers-color-scheme: light)">
<meta name="theme-color" content="#111317" media="(prefers-color-scheme: dark)">

<!-- 6. Kafel Windows (opcjonalnie — bez pliku browserconfig.xml) -->
<meta name="msapplication-TileColor" content="#111317">
<meta name="msapplication-TileImage" content="/icon-192.png">

<!-- 7. Nazwa aplikacji w trybie samodzielnym -->
<meta name="application-name" content="Danaco Console">
<meta name="apple-mobile-web-app-title" content="Danaco Console">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
"""
zapisz(os.path.join(KAT_FAVICON, "naglowek-snippet.html"), SNIPPET)
licznik["inne"] += 1

# Kontrola końcowa sprawdza, czy kropka sygnału przetrwała rasteryzację
# w najmniejszym rozmiarze, mierząc jej obecność w wytworzonych plikach.
kontrola = []
for nazwa in ("favicon-16.png", "favicon-32.png", "favicon-48.png"):
    im = Image.open(os.path.join(KAT_FAVICON, nazwa)).convert("RGB")
    niebieskie = sum(
        1 for p in im.getdata() if p[2] > p[0] + 30 and p[2] > 120
    )
    kontrola.append((nazwa, im.size[0], niebieskie))

print("Wytworzono:", licznik)
print("Kontrola kropki sygnału w rastrach faviconu (liczba pikseli błękitu):")
for nazwa, px, n in kontrola:
    print("  %-18s %3d px → %d px błękitu %s" % (nazwa, px, n, "OK" if n else "BRAK"))
