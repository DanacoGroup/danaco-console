/*! @license
 * Zestaw ikon Danaco Console wywodzi się z biblioteki Lucide (licencja ISC).
 * Copyright (c) 2026 Lucide Icons and Contributors
 * Portions: Feather — The MIT License (MIT), Copyright (c) 2013-present Cole Bemis
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 * Pełny tekst obu licencji: `ikony/LICENCJA-lucide.txt`.
 * Nota musi towarzyszyć każdej kopii — dlatego stoi w module wchodzącym
 * do pakietu wyjściowego, nie tylko w pliku obok źródeł.
 */
// Wiązanie plików SVG zestawu z nazwami ikon.
// Treść każdej ikony wczytywana jest jako tekst (`?raw`), ponieważ tylko SVG
// osadzony w dokumencie dziedziczy barwę tekstu przez `currentColor`
// i tym samym działa w obu motywach bez zmiany reguł.
// Wykaz nazw jest tu jedynym źródłem prawdy — typ `NazwaIkony` wywodzi się z niego.
// Same wiązania mieszkają w `zrodla/`, w podziale na grupy zastosowań
// zgodnym z porządkiem `manifest.json`; ten plik wyłącznie je scala.

import logoDanaco from './svg/logo-danaco.svg?raw';
import { CZYNNOSCI } from './zrodla/czynnosci';
import { NAWIGACJA } from './zrodla/nawigacja';
import { ROLE } from './zrodla/role';
import { SRODOWISKA } from './zrodla/srodowiska';
import { STANY } from './zrodla/stany';
import { ZASOBY } from './zrodla/zasoby';

/**
 * Zestaw ikon pakietu design v2.0 — 82 pozycje, dokładnie wykaz
 * `ikony/manifest.json`. Siatka 24×24, obrys 1.75, `currentColor`.
 */
export const ZESTAW_IKON = {
  ...NAWIGACJA,
  ...CZYNNOSCI,
  ...STANY,
  ...ZASOBY,
  ...ROLE,
  ...SRODOWISKA,
} as const;

/**
 * Godło marki — jedyny znak o barwach własnych, poza zestawem 82 ikon.
 * Plik `svg/logo-danaco.svg` to sygnet w odmianie na podłoże ciemne, bo oba
 * miejsca osadzenia (pasek górny, blok tożsamości) leżą na pasku
 * atramentowym w obu motywach. Pozostałe odmiany znaku — `ikony/marka.ts`.
 */
const GODLO = {
  'logo-danaco': logoDanaco,
} as const;

/**
 * Nazwy wycofane pakietem v2.0, nadal wołane przez widoki jeszcze
 * nieprzebudowane. Wskazują na następców z zestawu — żadnego nowego rysunku
 * ani pliku spoza pakietu. Znikają wraz z ostatnim wywołaniem.
 */
const ZGODNOSC_WSTECZNA = {
  'strzalka-dol': ZESTAW_IKON['grot-dol'],
  'karta-przegladarki': ZESTAW_IKON['karta-okna'],
  'waga': ZESTAW_IKON['debata'],
} as const;

/** Rejestr rozpoznawanych nazw: nazwa → źródło SVG. */
export const ZRODLA_IKON = {
  ...ZESTAW_IKON,
  ...GODLO,
  ...ZGODNOSC_WSTECZNA,
} as const;

/** Nazwa ikony należąca do zestawu. Literówka jest błędem kompilacji. */
export type NazwaIkony = keyof typeof ZRODLA_IKON;

/** Nazwa pozycji właściwego zestawu — bez godła i nazw wycofanych. */
export type NazwaZestawu = keyof typeof ZESTAW_IKON;

/** Nazwy wszystkich 82 ikon zestawu, w porządku manifestu. */
export const NAZWY_IKON = Object.keys(ZESTAW_IKON) as NazwaZestawu[];

/** Godło marki — jedyna pozycja pełnokolorowa. */
export const IKONA_PELNOKOLOROWA: NazwaIkony = 'logo-danaco';
