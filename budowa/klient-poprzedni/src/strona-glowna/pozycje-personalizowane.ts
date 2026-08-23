import { ComponentKind, type Component } from '../../../shared/contract';
import { ikonaRodzaju, type KodKomponentu, type PozycjaKomponentu } from './pozycje-komponentow';

/**
 * Kafle personalizowane strefy drugiej — komponenty już zbudowane, po jednym na
 * komponent. Kafel rodzaju z czwórki stałej prowadzi do zbudowania nowego;
 * kafel personalizowany wskazuje konkretną automatykę, eksperta czy projekt.
 *
 * Na kaflu stoi `Component.name` — nazwa nadana przy `component.create`. Klient
 * jej nie skraca, nie poprawia i nie podmienia na nazwę rodzaju.
 *
 * Kafel nie jest osobnym bytem, tylko widokiem komponentu: wykaz przychodzi
 * komendą `component.list` i nie ma osobnej komendy „utworzenia kafla".
 *
 * `component.list` bez `includeDisabled` oddaje wyłącznie komponenty czynne,
 * więc komponent wyłączony daje krótszą listę, a nie kafel wyszarzony.
 */

/** Wezwanie kafla, gdy komponent nie ma opisu własnego. */
const WEZWANIA: Record<KodKomponentu, string> = {
  automations: 'Otwórz automatykę',
  agents: 'Otwórz eksperta',
  workspace: 'Otwórz projekt',
  assistant: 'Otwórz profil asystenta',
};

/** Czy rodzaj z rdzenia należy do zamkniętego zbioru kontraktu. */
function rodzajZnany(rodzaj: string): rodzaj is KodKomponentu {
  return (Object.values(ComponentKind) as readonly string[]).includes(rodzaj);
}

/**
 * Komponent rdzenia jako kafel personalizowany.
 *
 * Komponent rodzaju nieznanego klientowi zwraca `null` — nie dlatego, że jest
 * zakazany, tylko dlatego, że kafel nie wiedziałby, dokąd prowadzi. Kafel bez
 * skutku byłby atrapą, a krótsza lista atrapą nie jest.
 */
export function pozycjaZKomponentu(komponent: Component): PozycjaKomponentu | null {
  if (!rodzajZnany(komponent.kind)) return null;
  const opis = komponent.description ?? '';
  return {
    kod: komponent.kind,
    nazwa: komponent.name,
    wezwanie: opis !== '' ? opis : WEZWANIA[komponent.kind],
    ikona: ikonaRodzaju(komponent.kind),
    // Metadane idą z odpowiedzi, nie z domysłu — `component.list` oddaje je
    // przy każdej pozycji.
    metadane: {
      czynny: komponent.enabled,
      zmieniony: komponent.updatedAt,
      ...(komponent.targetId === undefined ? {} : { cel: komponent.targetId }),
    },
    komponent: komponent.id,
  };
}

/** Wykaz komponentów rdzenia jako kafle, z pominięciem rodzajów nieznanych. */
export function pozycjeZKomponentow(
  komponenty: readonly Component[],
): readonly PozycjaKomponentu[] {
  const pozycje: PozycjaKomponentu[] = [];
  for (const komponent of komponenty) {
    const pozycja = pozycjaZKomponentu(komponent);
    if (pozycja !== null) pozycje.push(pozycja);
  }
  return pozycje;
}
