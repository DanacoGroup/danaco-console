import type { OknoKomunikacji } from '../okno-komunikacji/okno';
import { wpisModelu, type Wpis } from '../okno-komunikacji/wpis';

/** Zapis historii gniazda, prowadzony równolegle z widokiem okna, zapisujący wpisy i fragmenty wypowiedzi wraz z odtwarzaniem ich w nowym widoku. */
export interface DziennikWpisow {
  /** Odnotowuje wpis dołożony do historii. */
  zapisz(wpis: Wpis): void;
  /** Odnotowuje fragment strumienia — dokładany do wypowiedzi tej samej persony. */
  zapiszFragment(persona: string, fragment: string): void;
  /** Odtwarza całą historię w świeżo zbudowanym oknie. */
  odtworz(okno: OknoKomunikacji): void;
  /** Liczba odnotowanych wpisów. */
  liczba(): number;
}

/**
 * Dziennik wpisów gniazda trzyma zapis historii, żeby przebudowa widoku okna komunikacji po zmianie roli mogła odtworzyć rozmowę, zamiast ją skasować.
 */
export function utworzDziennikWpisow(): DziennikWpisow {
  const wpisy: Wpis[] = [];

  return {
    zapisz(wpis) {
      wpisy.push({ ...wpis });
    },

    zapiszFragment(persona, fragment) {
      const ostatni = wpisy.at(-1);
      if (ostatni === undefined || ostatni.persona !== persona) {
        wpisy.push(wpisModelu(fragment, persona));
        return;
      }
      ostatni.tresc = `${ostatni.tresc}${fragment}`;
    },

    odtworz(okno) {
      for (const wpis of wpisy) okno.dopisz({ ...wpis });
    },

    liczba: () => wpisy.length,
  };
}
