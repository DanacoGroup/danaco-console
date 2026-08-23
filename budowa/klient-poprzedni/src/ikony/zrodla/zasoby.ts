// Źródła SVG — grupa: zasoby.
// Nazwy i kolejność wprost z `ikony/manifest.json` (pozycje 41–54).

import dokument from '../svg/dokument.svg?raw';
import plik from '../svg/plik.svg?raw';
import folder from '../svg/folder.svg?raw';
import archiwum from '../svg/archiwum.svg?raw';
import obraz from '../svg/obraz.svg?raw';
import kod from '../svg/kod.svg?raw';
import koperta from '../svg/koperta.svg?raw';
import kalendarz from '../svg/kalendarz.svg?raw';
import tabelaDanych from '../svg/tabela-danych.svg?raw';
import biblioteka from '../svg/biblioteka.svg?raw';
import baza from '../svg/baza.svg?raw';
import wykres from '../svg/wykres.svg?raw';
import terminal from '../svg/terminal.svg?raw';
import galaz from '../svg/galaz.svg?raw';

/** Zasoby i dane — dokumenty, katalogi, zestawienia, repozytoria. */
export const ZASOBY = {
  'dokument': dokument,
  'plik': plik,
  'folder': folder,
  'archiwum': archiwum,
  'obraz': obraz,
  'kod': kod,
  'koperta': koperta,
  'kalendarz': kalendarz,
  'tabela-danych': tabelaDanych,
  'biblioteka': biblioteka,
  'baza': baza,
  'wykres': wykres,
  'terminal': terminal,
  'galaz': galaz,
} as const;
