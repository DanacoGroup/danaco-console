import {
  TelemetryFormat,
  type ProvenanceTraceExportRequest,
  type ProvenanceTraceExportResponse,
} from '../../../../shared/contract';
import {
  pobierzPlik,
  przelacznik,
  przyciskAkcji,
  wiersz,
  wybor,
} from '../../modele/kontrolki-formularza-braki';
import { zdanieNiepowodzenia } from './niepowodzenie-odczytu';
import {
  FORMAT_WYDANIA,
  RODZAJ_TRESCI_WYDANIA,
  ROZSZERZENIE_WYDANIA,
} from './prowenancja-slowa';
import type { ZrodloProwenancji } from './prowenancja-zrodlo';

/**
 * Interfejs wydania śladu wywołań mówi z góry, co plik poniesie, i po jego
 * powstaniu potwierdza to na podstawie tego, co rdzeń rzeczywiście oddał.
 */
export interface WydanieSladu {
  /** Kontrolki paska: format, treść i przycisk wydania. */
  element: HTMLElement;
  /** Zdanie o tym, co wydanie poniesie przy dzisiejszych nastawach. */
  zapowiedz(): string;
}

/** Zakres wydania podawany przez zakładkę: te same filtry wywołań i przedziału czasu, które złożyły się na bieżący wykaz. */
export interface ZakresWydania {
  /** Wywołania wskazane wprost; puste znaczy „cały zakres czasu”. */
  wywolania: readonly string[];
  od?: number;
  do?: number;
}

export function utworzWydanieSladu(
  zrodlo: ZrodloProwenancji,
  zakres: () => ZakresWydania,
  poWydaniu: (zdanie: string, udane: boolean) => void,
): WydanieSladu {
  const format = wybor(
    'Format wydania śladu',
    Object.values(TelemetryFormat).map((wartosc) => [wartosc, FORMAT_WYDANIA[wartosc]]),
  );
  format.value = TelemetryFormat.Otlp;

  const zTrescia = przelacznik('Wydaj wraz z treścią promptu i odpowiedzi');

  const wydaj = przyciskAkcji('Wydaj ślad wywołań do pliku', 'dn-btn dn-btn--atrament');

  function zapowiedz(): string {
    const wybrany = format.value as TelemetryFormat;
    const teraz = zakres();
    const czego =
      teraz.wywolania.length === 0
        ? 'wszystkie wywołania z zakresu czasu okna'
        : `wywołań wskazanych: ${String(teraz.wywolania.length)}`;
    return (
      `Wydanie obejmie ${czego}, w formacie: ${FORMAT_WYDANIA[wybrany]}. ` +
      (zTrescia.checked
        ? 'Plik poniesie TREŚĆ promptów i odpowiedzi — czyli materiał rozmów tej instalacji. ' +
          'O redakcji danych wrażliwych rozstrzyga ustawienie rdzenia, nie to okno; ' +
          'czy zaszła, powie zdanie po wydaniu.'
        : 'Plik poniesie sam ślad — kanał, model, czasy, tokeny i koszt — bez treści ' +
          'promptów i odpowiedzi.')
    );
  }

  const opis = document.createElement('p');
  opis.className = 'dn-pole-opis dg-narzedzie__opis';
  opis.dataset['zapowiedzWydania'] = 'tak';
  const odswiezOpis = (): void => void (opis.textContent = zapowiedz());
  odswiezOpis();
  format.addEventListener('change', odswiezOpis);
  zTrescia.addEventListener('change', odswiezOpis);

  wydaj.addEventListener('click', () => {
    const wybrany = format.value as TelemetryFormat;
    const teraz = zakres();
    const zadanie: ProvenanceTraceExportRequest = {
      format: wybrany,
      ...(teraz.wywolania.length === 0 ? {} : { callIds: [...teraz.wywolania] }),
      ...(teraz.od === undefined ? {} : { fromTime: teraz.od }),
      ...(teraz.do === undefined ? {} : { toTime: teraz.do }),
      ...(zTrescia.checked ? { includeContent: true } : {}),
    };
    poWydaniu(`Wysłano żądanie wydania śladu. ${zapowiedz()}`, true);
    void zrodlo.wydajSlad(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        poWydaniu(zdanieNiepowodzenia('wydania śladu wywołań', wynik.powod), false);
        return;
      }
      const oddane = wynik.wynik;
      if (oddane.callCount === 0) {
        // Zero wywołań w wydaniu nie jest odmową ani plikiem: pusty plik udawałby ślad, którego nie ma.
        poWydaniu(
          'Rdzeń przyjął żądanie i objął wydaniem ZERO wywołań — w tym zakresie nie ma czego ' +
            'wydać, więc plik nie powstał.',
          false,
        );
        return;
      }
      pobierzPlik(
        nazwaPliku(oddane.format),
        oddane.content,
        RODZAJ_TRESCI_WYDANIA[oddane.format],
      );
      poWydaniu(zdanieWydania(oddane, zTrescia.checked), true);
    });
  });

  const element = document.createElement('div');
  element.className = 'dg-narzedzie__pasek';
  element.append(
    wiersz('Format wydania śladu', format, {
      klasa: 'dg-wiersz',
      objasnienie:
        'Format jedzie w polu format komendy provenance.trace.export; wszystkie cztery ' +
        'wartości pochodzą z kontraktu.',
    }),
    wiersz('Wydaj wraz z treścią promptu i odpowiedzi', zTrescia, {
      klasa: 'dg-wiersz',
      objasnienie:
        'Pole includeContent. Kontrakt zakłada brak treści, dopóki Operator jej nie zażąda — ' +
        'przełącznik startuje wyłączony.',
    }),
    wydaj,
    opis,
  );

  return { element, zapowiedz };
}

/** Nazwa pliku wydania złożona z rozszerzenia formatu i znacznika czasu, tak aby sama nazwa mówiła, co plik niesie i z kiedy jest. */
function nazwaPliku(format: TelemetryFormat): string {
  const stempel = new Date().toISOString().replace(/[:.]/gu, '-');
  return `slad-wywolan-modelu-${stempel}.${ROZSZERZENIE_WYDANIA[format]}`;
}

/**
 * Zdanie o pliku, który właśnie powstał — z pól odpowiedzi, nie z zamówienia.
 *
 * Rdzeń ma prawo wydać w innym formacie niż poproszono i ma prawo zredagować
 * treść; jedno i drugie musi być widoczne, bo od tego zależy, czym ten plik
 * wolno się posłużyć.
 */
function zdanieWydania(oddane: ProvenanceTraceExportResponse, zamowionoTresc: boolean): string {
  const czesci = [
    `Rdzeń wydał ślad ${String(oddane.callCount)} wywołań w formacie ${FORMAT_WYDANIA[oddane.format]};` +
      ` znaków treści: ${String(oddane.content.length)}. Plik zszedł na urządzenie Operatora.`,
  ];
  if (zamowionoTresc) {
    czesci.push(
      oddane.redacted === true
        ? 'Treść promptów i odpowiedzi jest w pliku PO REDAKCJI danych wrażliwych.'
        : 'Treść promptów i odpowiedzi jest w pliku BEZ redakcji danych wrażliwych — rdzeń nie ' +
          'orzekł redakcji, więc plik niesie materiał rozmów w postaci pierwotnej.',
    );
  } else {
    czesci.push('Treści promptów i odpowiedzi plik nie niesie — nie była zamówiona.');
  }
  return czesci.join(' ');
}
