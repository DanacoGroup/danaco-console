#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# Generator sygnatury poczty elektronicznej modułu Danaco Console, wytwarzający
# pliki HTML w wersji jasnej i ciemnej.
"""
Danaco Console — generator sygnatury poczty elektronicznej (B6).

Buduje dwa gotowe do wklejenia pliki HTML (wersja jasna i ciemna). Znak jest
osadzony jako PNG w data URI (2×), bo klienci poczty blokują SVG i odnośniki
zewnętrzne. Układ oparty na tabeli, style wyłącznie w atrybucie `style` — to
jedyna droga, którą przechodzą Outlook (silnik Word), Gmail i Apple Mail.

Uruchomienie: python3 generator-sygnatury.py
"""

import base64
import importlib.util
import os

import cairosvg

BAZA = os.path.dirname(os.path.abspath(__file__))
KAT_HTML = os.path.join(BAZA, "html")
KAT_PNG = os.path.join(BAZA, "png")
os.makedirs(KAT_HTML, exist_ok=True)

_spec = importlib.util.spec_from_file_location(
    "gen_zast", os.path.join(BAZA, "generator-zastosowan.py"))
gz = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(gz)


def _lockup_png(ink, dot, szer_px):
    """Lockup kompaktowy jako PNG bez tła, o zadanej szerokości."""
    skala = 1.0
    svg = ('<svg xmlns="http://www.w3.org/2000/svg" width="%.2f" height="96" '
           'viewBox="0 0 %.2f 96">%s</svg>'
           % (gz.LOCKUP_KOMPAKT_SZER, gz.LOCKUP_KOMPAKT_SZER,
              gz.lockup_kompaktowy(0, 0, skala, ink, dot)))
    return cairosvg.svg2png(bytestring=svg.encode("utf-8"), output_width=szer_px)


def _data_uri(dane):
    return "data:image/png;base64," + base64.b64encode(dane).decode("ascii")


SZER_ZNAKU = 168          # szerokość znaku w sygnaturze (px CSS)


def _wysokosc_znaku():
    """Wysokość odczytana z faktycznego rastra 2x — atrybut `height` musi
    zgadzać się z plikiem, inaczej Outlook nieznacznie ściska znak."""
    import io
    from PIL import Image
    obraz = Image.open(io.BytesIO(_lockup_png(gz.INK_J, gz.DOT_J, SZER_ZNAKU * 2)))
    return int(round(obraz.size[1] / 2.0))


WYS_ZNAKU = _wysokosc_znaku()

SZABLON = """<!doctype html>
<html lang="pl">
<head>
<meta charset="utf-8">
<title>Danaco Console — sygnatura poczty elektronicznej ({wariant})</title>
</head>
<body style="margin:0;padding:24px;background:{tlo_strony};">
<!-- ══════════════════════════════════════════════════════════════════════
     DANACO CONSOLE — SYGNATURA POCZTY ELEKTRONICZNEJ · wariant {wariant}
     Skopiuj CAŁY blok <table> poniżej i wklej w ustawieniach klienta poczty.
     Pola w nawiasach kwadratowych uzupełnij własnymi danymi.
     Zasady: szerokość 100%% do 520 px · znak {szer}×{wys} px (2× w pliku) ·
     wiersz reguły 1 px · kropka sygnału jedynym akcentem barwnym.
     ══════════════════════════════════════════════════════════════════════ -->
<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="border-collapse:collapse;width:100%;max-width:520px;background:{tlo};font-family:'IBM Plex Sans','Segoe UI',Arial,sans-serif;">
  <tr>
    <td style="padding:16px 20px 12px 20px;">
      <img src="{znak}" width="{szer}" height="{wys}" alt="Danaco Console" style="display:block;border:0;outline:none;text-decoration:none;width:{szer}px;height:auto;">
    </td>
  </tr>
  <tr>
    <td style="padding:0 20px 4px 20px;font-family:'IBM Plex Sans','Segoe UI',Arial,sans-serif;font-size:15px;line-height:20px;font-weight:600;color:{tekst};">
      [Imię i nazwisko]
    </td>
  </tr>
  <tr>
    <td style="padding:0 20px 12px 20px;font-family:'IBM Plex Mono','Courier New',monospace;font-size:11px;line-height:16px;letter-spacing:1.4px;text-transform:uppercase;color:{tekst3};">
      [stanowisko]
    </td>
  </tr>
  <tr>
    <td style="padding:0 20px 0 20px;">
      <table role="presentation" cellpadding="0" cellspacing="0" border="0" style="border-collapse:collapse;width:100%;">
        <tr><td height="1" style="height:1px;line-height:1px;font-size:0;background:{obrys};">&nbsp;</td></tr>
      </table>
    </td>
  </tr>
  <tr>
    <td style="padding:12px 20px 0 20px;font-family:'IBM Plex Sans','Segoe UI',Arial,sans-serif;font-size:13px;line-height:19px;color:{tekst2};">
      {producent}<br>
      <a href="mailto:{kontakt}" style="color:{sygnal};text-decoration:none;">{kontakt}</a>
    </td>
  </tr>
  <tr>
    <td style="padding:10px 20px 18px 20px;font-family:'IBM Plex Mono','Courier New',monospace;font-size:10px;line-height:15px;letter-spacing:1.6px;text-transform:uppercase;color:{tekst3};">
      Danaco Console &middot; {deskryptor}
    </td>
  </tr>
</table>
<!-- ══════════════ KONIEC SYGNATURY ══════════════ -->
</body>
</html>
"""


def zbuduj(wariant, tlo, tlo_strony, tekst, tekst2, tekst3, obrys, sygnal, ink, dot):
    znak = _data_uri(_lockup_png(ink, dot, SZER_ZNAKU * 2))
    html = SZABLON.format(
        wariant=wariant, tlo=tlo, tlo_strony=tlo_strony, tekst=tekst,
        tekst2=tekst2, tekst3=tekst3, obrys=obrys, sygnal=sygnal, znak=znak,
        szer=SZER_ZNAKU, wys=WYS_ZNAKU, producent=gz.PRODUCENT,
        kontakt=gz.KONTAKT, deskryptor=gz.DESKRYPTOR)
    sciezka = os.path.join(KAT_HTML, "sygnatura-poczty-%s.html" % wariant)
    with open(sciezka, "w", encoding="utf-8") as plik:
        plik.write(html)
    return sciezka


def main():
    # wersja jasna — tło białe, atrament #181818
    zbuduj("jasna", tlo="#FFFFFF", tlo_strony="#F4F4F4", tekst="#181818",
           tekst2="#616161", tekst3="#7C7C7C", obrys="#E3E3E3",
           sygnal="#2457C9", ink=gz.INK_J, dot=gz.DOT_J)
    # wersja ciemna — tło #131313 (rama kokpitu), biel źródłowa
    zbuduj("ciemna", tlo="#131313", tlo_strony="#0F0F0F", tekst="#ECECEC",
           tekst2="#9E9E9E", tekst3="#7C7C7C", obrys="#2A2A2A",
           sygnal="#8FB2F5", ink=gz.INK_C, dot=gz.DOT_C)
    # znaki pomocnicze do ręcznego wklejenia
    for nazwa, ink, dot in (("znak-poczty-jasny", gz.INK_J, gz.DOT_J),
                            ("znak-poczty-ciemny", gz.INK_C, gz.DOT_C)):
        with open(os.path.join(KAT_PNG, nazwa + "@2x.png"), "wb") as plik:
            plik.write(_lockup_png(ink, dot, SZER_ZNAKU * 2))
    print("sygnatury:", os.listdir(KAT_HTML))


if __name__ == "__main__":
    main()
