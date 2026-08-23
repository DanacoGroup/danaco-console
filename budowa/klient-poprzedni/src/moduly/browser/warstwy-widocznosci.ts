/**
 * Warstwy widoczności modułu Browser — stan stopniowego ujawniania funkcji.
 *
 * Opracowanie modułu (rozdz. 3.1) dzieli jego wyposażenie na cztery warstwy:
 * pierwsza jest widoczna bez interakcji, druga otwiera się znacznikiem
 * kontekstowym, trzecia menu operacji, a czwarta skrótem klawiszowym albo
 * trybem administracyjnym. Ten byt trzyma, co jest w tej chwili odsłonięte;
 * pasek kontekstu (`pasek-kontekstu.ts`) jest wyłącznie jego widokiem.
 *
 * Ujawnianie nie jest bramą. Element zwinięty nie jest zablokowany — jest
 * schowany, a każda jego czynność zostaje osiągalna wyzwalaczem albo skrótem.
 * Dlatego skrót warstwy czwartej działa również wtedy, gdy tryb administracyjny
 * jest wyłączony: tryb decyduje o obecności wyzwalacza w pasku, nie o dostępie
 * do funkcji.
 *
 * Skrót podpina się do elementu modułu, nie do dokumentu: w powłoce stoi obok
 * siebie kilka modułów, a nasłuch założony na dokumencie odpowiadałby także na
 * klawisze wciśnięte w cudzym oknie. Skrót działa więc wtedy, gdy fokus stoi
 * wewnątrz modułu. Powłoka gospodarza może przechwycić kombinację przed
 * stroną — wtedy pozostaje wyzwalacz w pasku.
 */

/** Numer warstwy widoczności z rozdziału 3.1 opracowania modułu. */
export type NumerWarstwy = 1 | 2 | 3 | 4;

/** Skrót klawiszowy z załącznika opracowania: klawisz przy Ctrl/Cmd i Shift. */
export interface SkrotWarstwy {
  /** Litera klawisza, tak jak podaje ją `KeyboardEvent.key`, bez rozróżnienia wielkości. */
  klawisz: string;
  /** Czy skrót wymaga klawisza Shift. */
  zShift: boolean;
}

/** Okno albo panel modułu odsłaniany wyzwalaczem paska kontekstu. */
export interface RozszerzenieModulu {
  /** Kod okna operacyjnego albo nazwa panelu — klucz stanu odsłonięcia. */
  kod: string;
  /** Napis wyzwalacza w pasku kontekstu. */
  nazwa: string;
  warstwa: NumerWarstwy;
  /** Element chowany i odsłaniany wraz ze zmianą stanu. */
  element: HTMLElement;
  /** Skrót klawiszowy rozszerzenia; pominięty znaczy „rozszerzenie skrótu nie ma". */
  skrot?: SkrotWarstwy;
  /**
   * Liczba pozycji, które rozszerzenie zebrało — staje przy nazwie wyzwalacza.
   * Pominięta znaczy, że rozszerzenie niczego nie zlicza; zero jest wtedy
   * nieodróżnialne od braku licznika, więc pominięcie ma znaczenie.
   */
  licznik?(): number;
}

export interface WarstwyWidocznosci {
  /** Dopisuje rozszerzenie i od razu nanosi na nie stan odsłonięcia. */
  zarejestruj(rozszerzenie: RozszerzenieModulu): void;
  /** Rozszerzenia w kolejności rejestracji — do zbudowania paska kontekstu. */
  rozszerzenia(): readonly RozszerzenieModulu[];
  czyOdsloniete(kod: string): boolean;
  /** Odsłania albo zwija rozszerzenie i oddaje stan po zmianie. */
  przelacz(kod: string): boolean;
  ustaw(kod: string, odsloniete: boolean): void;
  /** Czy pasek kontekstu pokazuje wyzwalacze warstwy czwartej. */
  trybAdministracyjny(): boolean;
  ustawTrybAdministracyjny(wlaczony: boolean): void;
  /** Podpina skróty rozszerzeń do elementu modułu; oddaje odpięcie. */
  podepnijSkroty(cel: HTMLElement): () => void;
}

export function utworzWarstwyWidocznosci(oglos: () => void): WarstwyWidocznosci {
  const wykaz: RozszerzenieModulu[] = [];
  const odsloniete = new Set<string>();
  let trybAdministracyjny = false;

  /** Nanosi stan odsłonięcia na element rozszerzenia. */
  function nanies(rozszerzenie: RozszerzenieModulu): void {
    rozszerzenie.element.hidden = !odsloniete.has(rozszerzenie.kod);
  }

  function ustaw(kod: string, czy: boolean): void {
    if (czy) odsloniete.add(kod);
    else odsloniete.delete(kod);
    for (const rozszerzenie of wykaz) nanies(rozszerzenie);
    oglos();
  }

  /** Rozszerzenie, którego skrót odpowiada wciśniętej kombinacji. */
  function podSkrot(zdarzenie: KeyboardEvent): RozszerzenieModulu | null {
    if (!(zdarzenie.ctrlKey || zdarzenie.metaKey)) return null;
    const klawisz = zdarzenie.key.toLowerCase();
    return (
      wykaz.find(
        (rozszerzenie) =>
          rozszerzenie.skrot !== undefined &&
          rozszerzenie.skrot.klawisz.toLowerCase() === klawisz &&
          rozszerzenie.skrot.zShift === zdarzenie.shiftKey,
      ) ?? null
    );
  }

  return {
    zarejestruj(rozszerzenie) {
      wykaz.push(rozszerzenie);
      nanies(rozszerzenie);
    },

    rozszerzenia: () => wykaz,
    czyOdsloniete: (kod) => odsloniete.has(kod),

    przelacz(kod) {
      const teraz = !odsloniete.has(kod);
      ustaw(kod, teraz);
      return teraz;
    },

    ustaw,
    trybAdministracyjny: () => trybAdministracyjny,

    ustawTrybAdministracyjny(wlaczony) {
      if (wlaczony === trybAdministracyjny) return;
      trybAdministracyjny = wlaczony;
      oglos();
    },

    podepnijSkroty(cel) {
      const nasluch = (zdarzenie: KeyboardEvent): void => {
        const rozszerzenie = podSkrot(zdarzenie);
        if (rozszerzenie === null) return;
        zdarzenie.preventDefault();
        ustaw(rozszerzenie.kod, !odsloniete.has(rozszerzenie.kod));
      };
      cel.addEventListener('keydown', nasluch);
      return () => cel.removeEventListener('keydown', nasluch);
    },
  };
}
