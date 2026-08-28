// Źródła wektorowe grupy ról: znaczniki `svg` wczytywane jako tekst i podane
// dalej bez przetwarzania. Nazwy oraz kolejność pozycji pochodzą wprost
// z pliku `ikony/manifest.json`.

import agent from '../svg/agent.svg?raw';
import agenci from '../svg/agenci.svg?raw';
import uzytkownik from '../svg/uzytkownik.svg?raw';
import siec from '../svg/siec.svg?raw';
import wezly from '../svg/wezly.svg?raw';
import wykonawca from '../svg/wykonawca.svg?raw';
import walidator from '../svg/walidator.svg?raw';
import warstwy from '../svg/warstwy.svg?raw';
import automatyzacja from '../svg/automatyzacja.svg?raw';
import debata from '../svg/debata.svg?raw';
import tlumacz from '../svg/tlumacz.svg?raw';
import badanie from '../svg/badanie.svg?raw';
import aplikacje from '../svg/aplikacje.svg?raw';
import paleta from '../svg/paleta.svg?raw';
import rozmowa from '../svg/rozmowa.svg?raw';
import mikrofon from '../svg/mikrofon.svg?raw';
import diagnostyka from '../svg/diagnostyka.svg?raw';
import monitor from '../svg/monitor.svg?raw';
import telefon from '../svg/telefon.svg?raw';
import cpu from '../svg/cpu.svg?raw';
import polecenie from '../svg/polecenie.svg?raw';
import slonce from '../svg/slonce.svg?raw';
import ksiezyc from '../svg/ksiezyc.svg?raw';
import globus from '../svg/globus.svg?raw';

/**
 * Role, moduły i środowisko wykonania: odwzorowanie nazwy ikony na treść
 * pliku wektorowego. Klucze odpowiadają nazwom plików katalogu `ikony/svg`,
 * a zapis `as const` zawęża je do zbioru dosłownego.
 */
export const ROLE = {
  'agent': agent,
  'agenci': agenci,
  'uzytkownik': uzytkownik,
  'siec': siec,
  'wezly': wezly,
  'wykonawca': wykonawca,
  'walidator': walidator,
  'warstwy': warstwy,
  'automatyzacja': automatyzacja,
  'debata': debata,
  'tlumacz': tlumacz,
  'badanie': badanie,
  'aplikacje': aplikacje,
  'paleta': paleta,
  'rozmowa': rozmowa,
  'mikrofon': mikrofon,
  'diagnostyka': diagnostyka,
  'monitor': monitor,
  'telefon': telefon,
  'cpu': cpu,
  'polecenie': polecenie,
  'slonce': slonce,
  'ksiezyc': ksiezyc,
  'globus': globus,
} as const;
