#!/usr/bin/env python3
# Skrypt uruchamiany bezpośrednio, interpreter systemowy wskazany wprost w tym pierwszym wierszu pliku wykonywalnego.
"""
Danaco Console — generator nośników marki (zastosowania).

Wytwarza z ZATWIERDZONEJ geometrii znaku (zero modyfikacji krzywych sygnetu):
  * obraz Open Graph 1200x630,
  * baner LinkedIn 1584x396,
  * tapety 2560x1440 i 3840x2160 (motyw ciemny i jasny),
  * tła slajdu 1920x1080 (tytułowe i treściowe),
  * kartę tytułową dokumentacji 1600x900,
  * wizytówkę 90x50 mm (awers, rewers, wersje ze spadami i paserami),
  * ekran powitalny 1280x800 i 2560x1600,
  * awatar 400x400 (kwadrat i koło),
  * znaki wodne o obniżonym kryciu.

Typografia jest zamieniana na krzywe (fontTools) — pliki SVG nie zależą od
zainstalowanych krojów i rasteryzują się identycznie wszędzie.

Uruchomienie:  python3 generator-zastosowan.py
Zależności:    fontTools + brotli (odczyt woff2), cairosvg, Pillow
"""

import os

import cairosvg
from fontTools.pens.svgPathPen import SVGPathPen
from fontTools.pens.transformPen import TransformPen
from fontTools.ttLib import TTFont

BAZA = os.path.dirname(os.path.abspath(__file__))
KAT_SVG = os.path.join(BAZA, "svg")
KAT_PNG = os.path.join(BAZA, "png")
KAT_FONTY = os.path.abspath(os.path.join(BAZA, "..", "..", "zasoby", "fonty"))
for k in (KAT_SVG, KAT_PNG):
    os.makedirs(k, exist_ok=True)

# ── barwy marki ─────────────────────────────────────────────────────────────
# W plikach SVG barwy podajemy dosłownie (to barwy własne znaku), obok podany
# jest odpowiadający żeton systemu — w HTML obowiązuje wyłącznie var(--dn-*).
TLO_CIEMNE = "#0F0F0F"   # --dn-tlo (ciemny)        = szary-950
TLO_JASNE = "#F4F4F4"    # --dn-tlo (jasny)         = szary-50
MEDALION = "#131313"     # --dn-rama                = szary-925
POWIERZCHNIA_C = "#181818"  # --dn-powierzchnia (c) = szary-900
INK_C = "#ECECEC"        # --dn-tekst (ciemny)      = szary-100
INK_J = "#181818"        # --dn-tekst (jasny)       = szary-900
TEKST2_C = "#9E9E9E"     # --dn-tekst-2 (ciemny)    = szary-400
TEKST2_J = "#616161"     # --dn-tekst-2 (jasny)     = szary-600
TEKST3 = "#7C7C7C"       # --dn-tekst-3 (oba)       = szary-500
DOT_C = "#5C8CEC"        # --dn-kropka (ciemny)     = sygnal-400
DOT_J = "#3B6FE0"        # --dn-kropka (jasny)      = sygnal-500

PRODUKT = "Danaco Console"
DESKRYPTOR = "AI Operating Environment"
ZDANIE = ("Zarzadzaj cyfrowa organizacja."
          .replace("Zarzadzaj", "Zarządzaj")
          .replace("cyfrowa", "cyfrową")
          .replace("organizacja", "organizacją"))
HIERARCHIA = "Środowisko · Moduł · Okno operacyjne"
SRODOWISKA = "TalkIn · WorkSpace · CodeStudio · MultitaskingAI"
PRODUCENT = "Danaco Holding Group Sp. z o.o."
KONTAKT = "support@danaco-group.pl"
WERSJA = "v2.0 · status deweloperski"

GROT_1 = "M12 26 H24 L44 48.0 L24 70 H12 L32 48.0 Z"
GROT_2 = "M40 26 H52 L72 48.0 L52 70 H40 L60 48.0 Z"
KROPKA = (83, 63.5, 6.5)            # cx, cy, r w siatce 96x96 (zapis jak w pliku źródłowym)
GROT_UPR = "M18 18 H35 L64 48.0 L35 78 H18 L45 48.0 Z"
KROPKA_UPR = (79, 69, 9)
OS_OPTYCZNA_X = 50.75               # środek farby sygnetu (12 → 89,5)
OS_OPTYCZNA_Y = 48.0

RODZINY = {
    "sg700": ["space-grotesk-latin-700-normal.woff2",
              "space-grotesk-latin-ext-700-normal.woff2"],
    "sg500": ["space-grotesk-latin-500-normal.woff2",
              "space-grotesk-latin-ext-500-normal.woff2"],
    "mono400": ["ibm-plex-mono-latin-400-normal.woff2",
                "ibm-plex-mono-latin-ext-400-normal.woff2"],
    "mono500": ["ibm-plex-mono-latin-500-normal.woff2",
                "ibm-plex-mono-latin-ext-500-normal.woff2"],
    "sans400": ["ibm-plex-sans-latin-400-normal.woff2",
                "ibm-plex-sans-latin-ext-400-normal.woff2"],
    "sans600": ["ibm-plex-sans-latin-600-normal.woff2",
                "ibm-plex-sans-latin-ext-600-normal.woff2"],
}
_PAMIEC = {}


def _fonty(rodzina):
    if rodzina not in _PAMIEC:
        pozycje = []
        for nazwa in RODZINY[rodzina]:
            font = TTFont(os.path.join(KAT_FONTY, nazwa))
            pozycje.append((font, font.getBestCmap(), font.getGlyphSet(), font["hmtx"]))
        _PAMIEC[rodzina] = pozycje
    return _PAMIEC[rodzina]


def _glif(rodzina, znak):
    for font, cmap, gset, hmtx in _fonty(rodzina):
        if ord(znak) in cmap:
            nazwa = cmap[ord(znak)]
            return gset, hmtx[nazwa][0], nazwa
    raise KeyError("brak glifu %r w rodzinie %s" % (znak, rodzina))


def szerokosc(tekst, rodzina, rozmiar, tracking=0.0):
    """Szerokość napisu w px (bez światła za ostatnim znakiem)."""
    skala = rozmiar / 1000.0
    suma = 0.0
    for i, znak in enumerate(tekst):
        _, adv, _ = _glif(rodzina, znak)
        suma += adv * skala
        if i < len(tekst) - 1:
            suma += tracking * rozmiar
    return suma


def tekst(tekst_zrodlowy, rodzina, rozmiar, x, y, barwa,
          tracking=0.0, kotwica="start", krycie=None):
    """Napis zamieniony na krzywe. y = linia bazowa. Kotwica: start|middle|end."""
    skala = rozmiar / 1000.0
    szer = szerokosc(tekst_zrodlowy, rodzina, rozmiar, tracking)
    if kotwica == "middle":
        x -= szer / 2.0
    elif kotwica == "end":
        x -= szer
    kursor = x
    kawalki = []
    for znak in tekst_zrodlowy:
        gset, adv, nazwa = _glif(rodzina, znak)
        if znak != " ":
            pioro = SVGPathPen(gset)
            gset[nazwa].draw(TransformPen(pioro, (skala, 0, 0, -skala, kursor, y)))
            d = pioro.getCommands()
            if d:
                kawalki.append(d)
        kursor += adv * skala + tracking * rozmiar
    if not kawalki:
        return ""
    atr_krycie = ' fill-opacity="%s"' % krycie if krycie is not None else ""
    return '<path fill="%s"%s d="%s"/>' % (barwa, atr_krycie, " ".join(kawalki))


def sygnet(x, y, bok, ink, dot, uproszczony=False, krycie_grotow=None,
           krycie_kropki=None):
    """Sygnet wpisany w kwadrat `bok` × `bok` (siatka 96 × 96), lewy górny róg (x, y)."""
    s = bok / 96.0
    kg = ' fill-opacity="%s"' % krycie_grotow if krycie_grotow is not None else ""
    kk = ' fill-opacity="%s"' % krycie_kropki if krycie_kropki is not None else ""
    if uproszczony:
        cx, cy, r = KROPKA_UPR
        srodek = ('<path fill="%s"%s d="%s"/>'
                  '<circle cx="%s" cy="%s" r="%s" fill="%s"%s/>'
                  % (ink, kg, GROT_UPR, cx, cy, r, dot, kk))
    else:
        cx, cy, r = KROPKA
        srodek = ('<path fill="%s"%s d="%s"/><path fill="%s"%s d="%s"/>'
                  '<circle cx="%s" cy="%s" r="%s" fill="%s"%s/>'
                  % (ink, kg, GROT_1, ink, kg, GROT_2, cx, cy, r, dot, kk))
    return '<g transform="translate(%.4f,%.4f) scale(%.6f)">%s</g>' % (x, y, s, srodek)


def sygnet_wysrodkowany(cx, cy, szer_farby, ink, dot, **kw):
    """Sygnet wyśrodkowany osią optyczną (50,75 / 48) w punkcie (cx, cy).
    `szer_farby` = szerokość samej farby (12 → 89,5 = 77,5 j.)."""
    s = szer_farby / 77.5
    bok = 96.0 * s
    x = cx - OS_OPTYCZNA_X * s
    y = cy - OS_OPTYCZNA_Y * s
    return sygnet(x, y, bok, ink, dot, **kw)


# Logotyp — metryka zatwierdzona (logotyp.svg): DANACO 34 px Space Grotesk 700
# z tracking 0,012 em, CONSOLE 12,5 px IBM Plex Mono 500 z krokiem 11,25 px
# (tracking 0,30 em), kropka kończąca (82; linia bazowa CONSOLE − 4) r = 3,4.
LOGO_DANACO_ROZM = 34.0
LOGO_DANACO_TRACK = 0.012
LOGO_CONSOLE_ROZM = 12.5
LOGO_CONSOLE_TRACK = 0.30
LOGO_BAZOWA_1 = 34.0      # linia bazowa DANACO w układzie logotypu
LOGO_BAZOWA_2 = 50.0      # linia bazowa CONSOLE
LOGO_KROPKA = (82.0, 46.0, 3.4)
LOGO_WYS_WERSALIKA = 0.700 * LOGO_DANACO_ROZM       # 23,8
LOGO_SZER = szerokosc  # alias czytelności


def logotyp(x, y, skala, ink, dot):
    """Logotyp (DANACO + CONSOLE + kropka). (x, y) = lewy górny róg wersalika
    DANACO; `skala` = 1,0 dla metryki źródłowej (wersalik 23,8 j.)."""
    dx = x
    dy = y - (LOGO_BAZOWA_1 - LOGO_WYS_WERSALIKA) * skala
    cx, cy, r = LOGO_KROPKA
    srodek = (
        tekst("DANACO", "sg700", LOGO_DANACO_ROZM, 0, LOGO_BAZOWA_1, ink,
              tracking=LOGO_DANACO_TRACK)
        + tekst("CONSOLE", "mono500", LOGO_CONSOLE_ROZM, 0, LOGO_BAZOWA_2, ink,
                tracking=LOGO_CONSOLE_TRACK)
        + '<circle cx="%s" cy="%s" r="%s" fill="%s"/>' % (cx, cy, r, dot))
    return '<g transform="translate(%.4f,%.4f) scale(%.6f)">%s</g>' % (dx, dy, skala, srodek)


def logotyp_szerokosc(skala):
    return szerokosc("DANACO", "sg700", LOGO_DANACO_ROZM, LOGO_DANACO_TRACK) * skala


def logotyp_wysokosc(skala):
    """Od górnej linii wersalika DANACO do linii bazowej CONSOLE."""
    return (LOGO_BAZOWA_2 - (LOGO_BAZOWA_1 - LOGO_WYS_WERSALIKA)) * skala


def lockup_poziomy(x, y, skala, ink, dot):
    """Konfiguracja C — metryka pliku logo-poziomy.svg (225 × 96 przy skali 1)."""
    srodek = (sygnet(0, 16, 64, ink, dot)
              + logotyp(76, 51 - LOGO_WYS_WERSALIKA, 1.0, ink, dot))
    return '<g transform="translate(%.4f,%.4f) scale(%.6f)">%s</g>' % (x, y, skala, srodek)


def lockup_pionowy(x, y, skala, ink, dot):
    """Konfiguracja D — metryka pliku logo-pionowy.svg (159 × 172 przy skali 1)."""
    srodek = (sygnet(31.7, 0, 96, ink, dot)
              + logotyp(12.0, 136 - LOGO_WYS_WERSALIKA, 1.0, ink, dot)
              .replace('translate(12.0000,', 'translate(12.0000,'))
    # CONSOLE jest wyśrodkowany pod DANACO, więc składamy go osobno.
    srodek = (sygnet(31.7, 0, 96, ink, dot)
              + tekst("DANACO", "sg700", LOGO_DANACO_ROZM, 12.0, 136.0, ink,
                      tracking=LOGO_DANACO_TRACK)
              + tekst("CONSOLE", "mono500", LOGO_CONSOLE_ROZM, 42.2, 154.0, ink,
                      tracking=LOGO_CONSOLE_TRACK)
              + '<circle cx="124.2" cy="150" r="3.4" fill="%s"/>' % dot)
    return '<g transform="translate(%.4f,%.4f) scale(%.6f)">%s</g>' % (x, y, skala, srodek)


def lockup_kompaktowy(x, y, skala, ink, dot):
    """Konfiguracja E — sygnet + DANACO + kropka kończąca, bez CONSOLE."""
    bazowa = 59.0
    szer_d = szerokosc("DANACO", "sg700", LOGO_DANACO_ROZM, LOGO_DANACO_TRACK)
    srodek = (sygnet(0, 16, 64, ink, dot)
              + tekst("DANACO", "sg700", LOGO_DANACO_ROZM, 76, bazowa, ink,
                      tracking=LOGO_DANACO_TRACK)
              + '<circle cx="%.2f" cy="%.2f" r="3.4" fill="%s"/>'
              % (76 + szer_d + 7.0, bazowa - 4.0, dot))
    return '<g transform="translate(%.4f,%.4f) scale(%.6f)">%s</g>' % (x, y, skala, srodek)


LOCKUP_POZIOMY_SZER = 76 + szerokosc("DANACO", "sg700", LOGO_DANACO_ROZM,
                                     LOGO_DANACO_TRACK)
LOCKUP_KOMPAKT_SZER = LOCKUP_POZIOMY_SZER + 10.4


def siatka(szer, wys, modul, barwa, krycie=0.045):
    linie = []
    x = float(modul)
    while x < szer:
        linie.append("M%.2f 0V%.2f" % (x, wys))
        x += modul
    y = float(modul)
    while y < wys:
        linie.append("M0 %.2fH%.2f" % (y, szer))
        y += modul
    return ('<path stroke="%s" stroke-opacity="%s" stroke-width="1" fill="none" d="%s"/>'
            % (barwa, krycie, "".join(linie)))


def dokument(szer, wys, tresc, tlo, jednostka=""):
    return ('<svg xmlns="http://www.w3.org/2000/svg" width="%s%s" height="%s%s" '
            'viewBox="0 0 %s %s" role="img" aria-label="Danaco Console">'
            '<rect width="%s" height="%s" fill="%s"/>%s</svg>'
            % (szer, jednostka, wys, jednostka, szer, wys, szer, wys, tlo, tresc))


def zapisz(nazwa, tresc_svg, png_szer=None):
    sciezka = os.path.join(KAT_SVG, nazwa + ".svg")
    with open(sciezka, "w", encoding="utf-8") as plik:
        plik.write(tresc_svg + "\n")
    if png_szer:
        cairosvg.svg2png(url=sciezka,
                         write_to=os.path.join(KAT_PNG, nazwa + ".png"),
                         output_width=png_szer)
    return sciezka


def obraz_og():
    W, H, M = 1200, 630, 80
    t = [siatka(W, H, 48, INK_C, 0.04)]
    t.append(lockup_poziomy(M, 64, 0.62, INK_C, DOT_C))
    t.append(tekst("Zarządzaj cyfrową", "sg700", 64, M, 300, INK_C, tracking=-0.01))
    szer_lin2 = szerokosc("organizacją", "sg700", 64, -0.01)
    t.append(tekst("organizacją", "sg700", 64, M, 376, INK_C, tracking=-0.01))
    t.append('<circle cx="%.1f" cy="%.1f" r="9" fill="%s"/>'
             % (M + szer_lin2 + 14, 376 - 9, DOT_C))
    t.append(tekst(DESKRYPTOR.upper(), "mono400", 17, M, 436, TEKST3, tracking=0.14))
    t.append('<path stroke="%s" stroke-opacity="0.5" stroke-width="1" d="M%s 496H%s"/>'
             % ("#2A2A2A", M, W - M))
    t.append(tekst(SRODOWISKA, "mono400", 16, M, 552, TEKST2_C, tracking=0.06))
    t.append(tekst("danaco console · v2.0", "mono400", 16, W - M, 552, TEKST3,
                   tracking=0.06, kotwica="end"))
    zapisz("obraz-og-1200x630", dokument(W, H, "".join(t), TLO_CIEMNE), 1200)

    # wariant jasny
    t = [siatka(W, H, 48, INK_J, 0.05)]
    t.append(lockup_poziomy(M, 64, 0.62, INK_J, DOT_J))
    t.append(tekst("Zarządzaj cyfrową", "sg700", 64, M, 300, INK_J, tracking=-0.01))
    t.append(tekst("organizacją", "sg700", 64, M, 376, INK_J, tracking=-0.01))
    t.append('<circle cx="%.1f" cy="%.1f" r="9" fill="%s"/>'
             % (M + szer_lin2 + 14, 376 - 9, DOT_J))
    t.append(tekst(DESKRYPTOR.upper(), "mono400", 17, M, 436, TEKST2_J, tracking=0.14))
    t.append('<path stroke="%s" stroke-opacity="1" stroke-width="1" d="M%s 496H%s"/>'
             % ("#E3E3E3", M, W - M))
    t.append(tekst(SRODOWISKA, "mono400", 16, M, 552, TEKST2_J, tracking=0.06))
    t.append(tekst("danaco console · v2.0", "mono400", 16, W - M, 552, TEKST3,
                   tracking=0.06, kotwica="end"))
    zapisz("obraz-og-1200x630-jasny", dokument(W, H, "".join(t), TLO_JASNE), 1200)


def baner_linkedin():
    W, H, M = 1584, 396, 120
    t = [siatka(W, H, 44, INK_C, 0.04)]
    t.append(sygnet(M, 132, 132, INK_C, DOT_C))
    xt = M + 132 + 56
    t.append(tekst("DANACO CONSOLE", "sg700", 52, xt, 196, INK_C, tracking=0.012))
    szer = szerokosc("DANACO CONSOLE", "sg700", 52, 0.012)
    t.append('<circle cx="%.1f" cy="%.1f" r="5.2" fill="%s"/>'
             % (xt + szer + 11, 190, DOT_C))
    t.append(tekst(DESKRYPTOR.upper(), "mono400", 17, xt, 242, TEKST2_C, tracking=0.14))
    t.append(tekst(HIERARCHIA, "mono400", 15, xt, 276, TEKST3, tracking=0.08))
    t.append(tekst(PRODUCENT, "mono400", 15, W - M, 242, TEKST2_C, tracking=0.04,
                   kotwica="end"))
    t.append(tekst(KONTAKT, "mono400", 15, W - M, 276, TEKST3, tracking=0.04,
                   kotwica="end"))
    zapisz("baner-linkedin-1584x396", dokument(W, H, "".join(t), TLO_CIEMNE), 1584)

    # wariant jasny
    t = [siatka(W, H, 44, INK_J, 0.05)]
    t.append(sygnet(M, 132, 132, INK_J, DOT_J))
    t.append(tekst("DANACO CONSOLE", "sg700", 52, xt, 196, INK_J, tracking=0.012))
    t.append('<circle cx="%.1f" cy="%.1f" r="5.2" fill="%s"/>'
             % (xt + szer + 11, 190, DOT_J))
    t.append(tekst(DESKRYPTOR.upper(), "mono400", 17, xt, 242, TEKST2_J, tracking=0.14))
    t.append(tekst(HIERARCHIA, "mono400", 15, xt, 276, TEKST3, tracking=0.08))
    t.append(tekst(PRODUCENT, "mono400", 15, W - M, 242, TEKST2_C, tracking=0.04,
                   kotwica="end"))
    t.append(tekst(KONTAKT, "mono400", 15, W - M, 276, TEKST3, tracking=0.04,
                   kotwica="end"))
    zapisz("baner-linkedin-1584x396-jasny", dokument(W, H, "".join(t), TLO_JASNE), 1584)


def tapeta(szer, wys, modul, ciemna=True):
    ink = INK_C if ciemna else INK_J
    dot = DOT_C if ciemna else DOT_J
    tlo = TLO_CIEMNE if ciemna else TLO_JASNE
    tekst_m = TEKST3
    kry_grot = 0.07 if ciemna else 0.09
    t = [siatka(szer, wys, modul, ink, 0.030 if ciemna else 0.045)]
    t.append(sygnet_wysrodkowany(szer * 0.5, wys * 0.46, szer * 0.215, ink, dot,
                                 krycie_grotow=kry_grot, krycie_kropki=0.90))
    podpis = "%s · %s" % (PRODUKT.upper(), DESKRYPTOR.upper())
    t.append(tekst(podpis, "mono400", max(12, szer / 150.0), szer * 0.5,
                   wys * 0.46 + szer * 0.108, tekst_m,
                   tracking=0.18, kotwica="middle", krycie=0.55))
    nazwa = "tapeta-%dx%d%s" % (szer, wys, "" if ciemna else "-jasna")
    zapisz(nazwa, dokument(szer, wys, "".join(t), tlo), szer)


def tla_slajdu():
    W, H, M = 1920, 1080, 120
    # 4a — slajd tytułowy; prawe dwie trzecie pozostają puste pod tytuł.
    t = [siatka(W, H, 60, INK_C, 0.035)]
    sk = 340.0 / 159.0
    t.append(lockup_pionowy(M, H / 2 - (172 * sk) / 2 - 40, sk, INK_C, DOT_C))
    t.append('<path stroke="#2A2A2A" stroke-width="1" d="M%s 760H%s"/>' % (M, M + 340))
    t.append(tekst(DESKRYPTOR.upper(), "mono400", 18, M, 808, TEKST3, tracking=0.14))
    t.append(tekst(PRODUCENT, "mono400", 15, M, H - 84, TEKST3, tracking=0.04))
    zapisz("tlo-slajdu-tytulowy-1920x1080", dokument(W, H, "".join(t), TLO_CIEMNE), 1920)

    # 4a′ — slajd tytułowy jasny
    t = [siatka(W, H, 60, INK_J, 0.05)]
    t.append(lockup_pionowy(M, H / 2 - (172 * sk) / 2 - 40, sk, INK_J, DOT_J))
    t.append('<path stroke="#E3E3E3" stroke-width="1" d="M%s 760H%s"/>' % (M, M + 340))
    t.append(tekst(DESKRYPTOR.upper(), "mono400", 18, M, 808, TEKST2_J, tracking=0.14))
    t.append(tekst(PRODUCENT, "mono400", 15, M, H - 84, TEKST3, tracking=0.04))
    zapisz("tlo-slajdu-tytulowy-1920x1080-jasny",
           dokument(W, H, "".join(t), TLO_JASNE), 1920)

    # 4b — slajd treściowy (pola tytułu i treści pozostają puste)
    t = [siatka(W, H, 60, INK_C, 0.025)]
    t.append('<path stroke="#212121" stroke-width="1" d="M%s 216H%s"/>' % (M, W - M))
    t.append(lockup_kompaktowy(W - M - 150, H - 96, 150 / LOCKUP_KOMPAKT_SZER,
                               INK_C, DOT_C))
    t.append(tekst(PRODUKT.upper(), "mono400", 14, M, H - 62, TEKST3, tracking=0.16))
    zapisz("tlo-slajdu-tresciowy-1920x1080", dokument(W, H, "".join(t), TLO_CIEMNE), 1920)

    # 4c — slajd treściowy jasny
    t = [siatka(W, H, 60, INK_J, 0.05)]
    t.append('<path stroke="#E3E3E3" stroke-width="1" d="M%s 216H%s"/>' % (M, W - M))
    t.append(lockup_kompaktowy(W - M - 150, H - 96, 150 / LOCKUP_KOMPAKT_SZER,
                               INK_J, DOT_J))
    t.append(tekst(PRODUKT.upper(), "mono400", 14, M, H - 62, TEKST3, tracking=0.16))
    zapisz("tlo-slajdu-tresciowy-1920x1080-jasny",
           dokument(W, H, "".join(t), TLO_JASNE), 1920)


def karta_tytulowa():
    W, H, M = 1600, 900, 100
    t = [siatka(W, H, 50, INK_C, 0.035)]
    t.append(lockup_poziomy(M, 84, 0.86, INK_C, DOT_C))
    t.append(tekst("Danaco Console", "sg700", 84, M, 452, INK_C, tracking=-0.01))
    t.append(tekst(DESKRYPTOR.upper(), "mono400", 20, M, 506, TEKST2_C, tracking=0.14))
    t.append('<path stroke="#2A2A2A" stroke-width="1" d="M%s 596H%s"/>' % (M, W - M))
    kolumny = [("WERSJA", "2.0"), ("STATUS", "deweloperski"),
               ("DATA", "2026-08-14"), ("PRODUCENT", "Danaco Holding Group")]
    x = M
    for etykieta, wartosc in kolumny:
        t.append(tekst(etykieta, "mono400", 13, x, 650, TEKST3, tracking=0.16))
        t.append(tekst(wartosc, "sans600", 19, x, 682, INK_C))
        x += 300
    t.append(tekst(HIERARCHIA, "mono400", 15, M, H - 72, TEKST3, tracking=0.08))
    zapisz("karta-tytulowa-dokumentacji-1600x900",
           dokument(W, H, "".join(t), TLO_CIEMNE), 1600)

    t = [siatka(W, H, 50, INK_J, 0.05)]
    t.append(lockup_poziomy(M, 84, 0.86, INK_J, DOT_J))
    t.append(tekst("Danaco Console", "sg700", 84, M, 452, INK_J, tracking=-0.01))
    t.append(tekst(DESKRYPTOR.upper(), "mono400", 20, M, 506, TEKST2_J, tracking=0.14))
    t.append('<path stroke="#E3E3E3" stroke-width="1" d="M%s 596H%s"/>' % (M, W - M))
    x = M
    for etykieta, wartosc in kolumny:
        t.append(tekst(etykieta, "mono400", 13, x, 650, TEKST3, tracking=0.16))
        t.append(tekst(wartosc, "sans600", 19, x, 682, INK_J))
        x += 300
    t.append(tekst(HIERARCHIA, "mono400", 15, M, H - 72, TEKST3, tracking=0.08))
    zapisz("karta-tytulowa-dokumentacji-1600x900-jasna",
           dokument(W, H, "".join(t), TLO_JASNE), 1600)


def wizytowka():
    W, H = 90.0, 50.0          # mm — układ współrzędnych = milimetry
    Mrg = 7.0
    # awers — atrament, lockup kompaktowy 28 mm
    skala = 28.0 / LOCKUP_KOMPAKT_SZER
    # Środek farby lockupu leży na y = 41,5, nie 48 — oś optyczna karty na 45% wysokości.
    t = [lockup_kompaktowy(Mrg, H * 0.45 - 41.5 * skala, skala, INK_C, DOT_C)]
    t.append(tekst(DESKRYPTOR.upper(), "mono400", 2.1, Mrg, H - Mrg, TEKST3,
                   tracking=0.14))
    zapisz("wizytowka-awers-90x50", dokument(W, H, "".join(t), MEDALION, "mm"), 1063)

    # rewers — jasny, blok kontaktowy
    t = [tekst("[Imię i nazwisko]", "sans600", 3.6, Mrg, 16.5, INK_J)]
    t.append(tekst("[stanowisko]", "mono400", 2.3, Mrg, 21.0, TEKST2_J, tracking=0.12))
    t.append('<path stroke="#E3E3E3" stroke-width="0.2" d="M%s 26H%s"/>'
             % (Mrg, W - Mrg))
    t.append(tekst(KONTAKT, "sans400", 2.8, Mrg, 32.5, INK_J))
    t.append(tekst(PRODUCENT, "sans400", 2.8, Mrg, 37.0, TEKST2_J))
    t.append(tekst("%s · %s" % (PRODUKT, DESKRYPTOR), "mono400", 2.1, Mrg,
                   H - Mrg, TEKST3, tracking=0.10))
    t.append(sygnet(W - Mrg - 11.0, Mrg, 11.0, INK_J, DOT_J))
    zapisz("wizytowka-rewers-90x50", dokument(W, H, "".join(t), TLO_JASNE, "mm"), 1063)

    # wersje ze spadem 3 mm i paserami
    for strona, tlo, ink in (("awers", MEDALION, INK_C), ("rewers", TLO_JASNE, INK_J)):
        S = 3.0
        Ws, Hs = W + 2 * S, H + 2 * S
        srodek = ('<g transform="translate(%s,%s)">%s</g>'
                  % (S, S, _wizytowka_srodek(strona)))
        pasery = []
        for (px, py, dx, dy) in ((S, 0, 0, 1), (S, Hs, 0, -1), (W + S, 0, 0, 1),
                                 (W + S, Hs, 0, -1), (0, S, 1, 0), (Ws, S, -1, 0),
                                 (0, H + S, 1, 0), (Ws, H + S, -1, 0)):
            pasery.append('<path stroke="%s" stroke-opacity="0.55" stroke-width="0.15" '
                          'd="M%s %sL%s %s"/>'
                          % (ink, px, py, px + dx * 2.0, py + dy * 2.0))
        zapisz("wizytowka-%s-90x50-spady" % strona,
               dokument(Ws, Hs, srodek + "".join(pasery), tlo, "mm"), 1134)


def _wizytowka_srodek(strona):
    W, H, Mrg = 90.0, 50.0, 7.0
    if strona == "awers":
        skala = 28.0 / LOCKUP_KOMPAKT_SZER
        return (lockup_kompaktowy(Mrg, H * 0.45 - 41.5 * skala, skala, INK_C, DOT_C)
                + tekst(DESKRYPTOR.upper(), "mono400", 2.1, Mrg, H - Mrg, TEKST3,
                        tracking=0.14))
    return (tekst("[Imię i nazwisko]", "sans600", 3.6, Mrg, 16.5, INK_J)
            + tekst("[stanowisko]", "mono400", 2.3, Mrg, 21.0, TEKST2_J, tracking=0.12)
            + '<path stroke="#E3E3E3" stroke-width="0.2" d="M%s 26H%s"/>' % (Mrg, W - Mrg)
            + tekst(KONTAKT, "sans400", 2.8, Mrg, 32.5, INK_J)
            + tekst(PRODUCENT, "sans400", 2.8, Mrg, 37.0, TEKST2_J)
            + tekst("%s · %s" % (PRODUKT, DESKRYPTOR), "mono400", 2.1, Mrg,
                    H - Mrg, TEKST3, tracking=0.10)
            + sygnet(W - Mrg - 11.0, Mrg, 11.0, INK_J, DOT_J))


def ekran_powitalny(szer, wys, modul):
    k = szer / 1280.0
    t = [siatka(szer, wys, modul, INK_C, 0.030)]
    skala_lock = (180.0 * k) / 159.0
    t.append(lockup_pionowy(szer / 2 - (159 * skala_lock) / 2,
                            wys * 0.5 - (172 * skala_lock) / 2 - 24 * k,
                            skala_lock, INK_C, DOT_C))
    # tor postępu — neutralny; jedynym akcentem pozostaje kropka znaku
    tor_w = 220.0 * k
    tor_y = wys * 0.5 + 128 * k
    t.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="%.1f" '
             'fill="%s" fill-opacity="0.10"/>'
             % (szer / 2 - tor_w / 2, tor_y, tor_w, 2 * k, 1 * k, INK_C))
    t.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="%.1f" '
             'fill="%s" fill-opacity="0.55"/>'
             % (szer / 2 - tor_w / 2, tor_y, tor_w * 0.42, 2 * k, 1 * k, INK_C))
    t.append(tekst("ŁĄCZENIE", "mono400", 13 * k, szer / 2, tor_y + 34 * k, TEKST3,
                   tracking=0.20, kotwica="middle"))
    t.append(tekst(WERSJA, "mono400", 12 * k, szer / 2, wys - 40 * k, TEKST3,
                   tracking=0.08, kotwica="middle"))
    zapisz("ekran-powitalny-%dx%d" % (szer, wys),
           dokument(szer, wys, "".join(t), TLO_CIEMNE), szer)


def awatar():
    R = 400
    znak = sygnet_wysrodkowany(R / 2, R / 2, R * 0.52, INK_C, DOT_C)
    zapisz("awatar-400x400-kwadrat", dokument(R, R, znak, MEDALION), 400)

    znak_kolo = sygnet_wysrodkowany(R / 2, R / 2, R * 0.44, INK_C, DOT_C)
    kolo = ('<svg xmlns="http://www.w3.org/2000/svg" width="400" height="400" '
            'viewBox="0 0 400 400" role="img" aria-label="Danaco Console">'
            '<circle cx="200" cy="200" r="200" fill="%s"/>%s</svg>' % (MEDALION, znak_kolo))
    zapisz("awatar-400x400-kolo", kolo, 400)

    # wariant jasny (koło na bieli — dla tła ciemnego serwisu)
    znak_j = sygnet_wysrodkowany(R / 2, R / 2, R * 0.44, INK_J, DOT_J)
    kolo_j = ('<svg xmlns="http://www.w3.org/2000/svg" width="400" height="400" '
              'viewBox="0 0 400 400" role="img" aria-label="Danaco Console">'
              '<circle cx="200" cy="200" r="200" fill="%s"/>%s</svg>'
              % (TLO_JASNE, znak_j))
    zapisz("awatar-400x400-kolo-jasny", kolo_j, 400)

    # wariant uproszczony (miniatura ≤ 48 px)
    znak_u = ('<g transform="translate(%.2f,%.2f) scale(%.5f)">'
              '<path fill="%s" d="%s"/><circle cx="%s" cy="%s" r="%s" fill="%s"/></g>'
              % (200 - 50.0 * 2.6667, 200 - 48.0 * 2.6667, 2.6667,
                 INK_C, GROT_UPR, KROPKA_UPR[0], KROPKA_UPR[1], KROPKA_UPR[2], DOT_C))
    zapisz("awatar-400x400-kwadrat-uproszczony",
           dokument(R, R, znak_u, MEDALION), 400)


def znaki_wodne():
    # sygnet — na jasnym (atrament 8 %) i na ciemnym (biel 10 %)
    for nazwa, ink, kry in (("znak-wodny-sygnet-na-jasnym", INK_J, 0.08),
                            ("znak-wodny-sygnet-na-ciemnym", INK_C, 0.10)):
        tresc = ('<svg xmlns="http://www.w3.org/2000/svg" width="96" height="96" '
                 'viewBox="0 0 96 96" role="img" aria-label="Danaco Console — znak wodny">'
                 '<g fill="%s" fill-opacity="%s">'
                 '<path d="%s"/><path d="%s"/><circle cx="%s" cy="%s" r="%s"/>'
                 '</g></svg>' % (ink, kry, GROT_1, GROT_2, *KROPKA))
        zapisz(nazwa, tresc, 512)

    # lockup poziomy — znak wodny dokumentu
    for nazwa, ink, kry in (("znak-wodny-lockup-na-jasnym", INK_J, 0.08),
                            ("znak-wodny-lockup-na-ciemnym", INK_C, 0.10)):
        srodek = lockup_poziomy(0, 0, 1.0, ink, ink)
        tresc = ('<svg xmlns="http://www.w3.org/2000/svg" width="225" height="96" '
                 'viewBox="0 0 225 96" role="img" aria-label="Danaco Console — znak wodny">'
                 '<g fill-opacity="%s">%s</g></svg>' % (kry, srodek))
        zapisz(nazwa, tresc, 900)

    # kafel powtarzalny do tła dokumentu (moduł 240 × 240, znak obrócony 0°)
    for nazwa, ink, kry in (("znak-wodny-kafel-na-jasnym", INK_J, 0.06),
                            ("znak-wodny-kafel-na-ciemnym", INK_C, 0.08)):
        srodek = sygnet_wysrodkowany(120, 120, 96, ink, ink)
        tresc = ('<svg xmlns="http://www.w3.org/2000/svg" width="240" height="240" '
                 'viewBox="0 0 240 240" role="img" aria-label="Danaco Console — kafel znaku wodnego">'
                 '<g fill-opacity="%s">%s</g></svg>' % (kry, srodek))
        zapisz(nazwa, tresc, 480)


def main():
    obraz_og()
    baner_linkedin()
    tapeta(2560, 1440, 80, True)
    tapeta(2560, 1440, 80, False)
    tapeta(3840, 2160, 120, True)
    tapeta(3840, 2160, 120, False)
    tla_slajdu()
    karta_tytulowa()
    wizytowka()
    ekran_powitalny(1280, 800, 40)
    ekran_powitalny(2560, 1600, 80)
    awatar()
    znaki_wodne()
    print("SVG :", len([p for p in os.listdir(KAT_SVG) if p.endswith(".svg")]))
    print("PNG :", len([p for p in os.listdir(KAT_PNG) if p.endswith(".png")]))


if __name__ == "__main__":
    main()
