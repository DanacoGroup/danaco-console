import { ComponentKind, type Component } from '../../../shared/contract';
import { ikonaRodzaju, type KodKomponentu, type PozycjaKomponentu } from './pozycje-komponentow';

/** Kafle personalizowane strefy drugiej pokazują komponenty już zbudowane, po jednym na komponent, z nazwą nadaną przy założeniu, której klient nie skraca ani nie podmienia. */
const WEZWANIA: Record<KodKomponentu, string> = {
  automations: 'Otwórz automatykę',
  agents: 'Otwórz eksperta',
  workspace: 'Otwórz projekt',
  assistant: 'Otwórz profil asystenta',
};

/** Rozstrzyga, czy rodzaj z rdzenia należy do zamkniętego zbioru kontraktu, wracając wtedy wartość prawdy. */
function rodzajZnany(rodzaj: string): rodzaj is KodKomponentu {
  return (Object.values(ComponentKind) as readonly string[]).includes(rodzaj);
}

/** Buduje kafel personalizowany z komponentu rdzenia; rodzaj nieznany klientowi zwraca brak, bo kafel nie wiedziałby, dokąd prowadzi. */
export function pozycjaZKomponentu(komponent: Component): PozycjaKomponentu | null {
  if (!rodzajZnany(komponent.kind)) return null;
  const opis = komponent.description ?? '';
  return {
    kod: komponent.kind,
    nazwa: komponent.name,
    wezwanie: opis !== '' ? opis : WEZWANIA[komponent.kind],
    ikona: ikonaRodzaju(komponent.kind),
    // Metadane idą z odpowiedzi, nie z domysłu.
    metadane: {
      czynny: komponent.enabled,
      zmieniony: komponent.updatedAt,
      ...(komponent.targetId === undefined ? {} : { cel: komponent.targetId }),
    },
    komponent: komponent.id,
  };
}

/** Przekłada wykaz komponentów rdzenia na kafle personalizowane, z pominięciem rodzajów nieznanych klientowi. */
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
