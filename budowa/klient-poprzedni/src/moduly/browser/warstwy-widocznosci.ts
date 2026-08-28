/**
 * Warstwy widoczności modułu Browser trzymają stan stopniowego ujawniania jego
 * wyposażenia: co jest w tej chwili odsłonięte i czy tryb administracyjny
 * pokazuje wyzwalacze warstwy czwartej. Pasek kontekstu jest wyłącznie widokiem
 * tego stanu.
 */

/**
 * Numer warstwy widoczności: pierwsza jest widoczna bez interakcji, druga otwiera
 * się znacznikiem kontekstowym, trzecia menu operacji, a czwarta skrótem
 * klawiszowym albo trybem administracyjnym.
 */
export type NumerWarstwy = 1 | 2 | 3 | 4;

/**
 * Skrót klawiszowy odsłaniający rozszerzenie: litera wciśnięta przy klawiszu
 * sterującym albo poleceniowym, opcjonalnie z klawiszem Shift. Klawisz sterujący
 * jest wymagany zawsze, żeby skrót nie kolidował z pisaniem w polu.
 */
export interface SkrotWarstwy {
  /** Litera klawisza, tak jak podaje ją `KeyboardEvent.key`, bez rozróżnienia wielkości. */
  klawisz: string;
  /** Czy skrót wymaga klawisza Shift. */
  zShift: boolean;
}

/**
 * Okno albo panel modułu odsłaniany wyzwalaczem paska kontekstu: niesie kod
 * będący kluczem stanu, napis wyzwalacza, numer warstwy, chowany element oraz
 * skrót klawiszowy i licznik zebranych pozycji, o ile je ma.
 */
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
  /** Liczba pozycji zebranych przez rozszerzenie; staje przy nazwie wyzwalacza. */
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
