// Panel wciąga swoje arkusze sam, bo gospodarzem bywa moduł, którego arkusz
// o Roundtable nic nie wie. `roundtable.css` daje trzy stany obowiązkowe
// nośnika `stany-okna`, `panel-debaty.css` — wygląd głosów. Arkusza `debata.css`
// panel nie wciąga: Debate Panel i Consensus Panel u gospodarza nie stoją.
import './roundtable.css';
import './panel-debaty.css';

import { Command } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import { przyciskAkcji, przyciskBezKomendy } from '../../modele/kontrolki-formularza';
import type { OpcjePanelu, PanelPomocniczy } from '../../okna-pomocnicze/panel-pomocniczy';
import { utworzRejestrKanalow } from '../../sterowanie/rejestr-kanalow';
import { powodBrakuObslugi } from './braki-kontraktu';
import { utworzStanDebaty, type StanDebaty } from './stan-debaty';
import { utworzStanTresci } from './stany-okna';
import { utworzStrumienWypowiedzi, type StrumienWypowiedzi } from './strumien-wypowiedzi';
import { rysujPanelDebaty } from './widok-panelu-debaty';
import { utworzZrodloRoundtable } from './zrodlo-roundtable';
import { utworzZrodloStrumieniaDebaty } from './zrodlo-strumienia-debaty';

/**
 * Panel debaty — przebieg debaty obok rozmowy, jako zwykły panel stosu.
 *
 * Kształt jest dokładnie ten, którego wymaga
 * `okna-pomocnicze/panel-pomocniczy.ts` — `element`, `odswiez()`, `zamknij()` —
 * i ani jedno pole ponad to. Gospodarzem bywa pas okien pomocniczych modułu
 * albo kolumna paneli sceny okien równoległych; panel żadnego z nich nie zna
 * i niczego o nich nie zakłada.
 *
 * Zamknięcie zdejmuje trzy rzeczy, nie jedną. Panel zakłada subskrypcję
 * `stream.chunk` (fragmenty wypowiedzi), subskrypcję zmian stanu debaty oraz —
 * przez `StanDebaty` — nasłuch `roundtable.debate.changed` i rejestru kanałów.
 * Wszystkie schodzą w `zamknij()`; panel bez tego zostawiłby je żywe po zejściu
 * ze sceny i rysowałby do elementu, którego nikt już nie ogląda.
 *
 * Własny stan debaty bierze się stąd, że `OpcjePanelu` daje wyłącznie `Kanal`,
 * `okno`, `modul` i `przedrostek` — egzemplarza `StanDebaty` tą drogą podać się
 * nie da, tak samo jak umowa `OpisModulu` nie przenosi rejestru kanałów. Drugi
 * stan nie jest drugą prawdą o debacie: czyta te same zdarzenia rdzenia i nie
 * ma ani jednej drogi zapisu — różni się od stanu złożenia wyłącznie chwilą
 * otwarcia nasłuchu.
 *
 * Panel nie odczytuje przebiegu na żądanie, choć kontrakt to przewiduje.
 * `roundtable.debate.get` oddaje skład, tury i wypowiedzi jednym wywołaniem,
 * ale obsługi tego odczytu jeszcze nie zbudowano — ani tutaj, ani w oknach
 * złożenia modułu — więc panel pokazuje wyłącznie to, co usłyszał od swojego
 * otwarcia. `odswiez()` odnawia zatem wykaz kanałów (nazwy uczestników)
 * i przerysowuje widok, a przycisk nazywający brakującą obsługę stoi widoczny
 * i klikalny zamiast zniknąć.
 */

/** Kod pozycji panelu — ten sam, którym rejestr i wytwórnia paneli go zawołają. */
export const KOD_PANELU_DEBATY = 'przebieg-debaty';

export function utworzPanelDebaty(opcje: OpcjePanelu): PanelPomocniczy {
  const rejestr = utworzRejestrKanalow(opcje.kanal);
  const zrodlo = utworzZrodloRoundtable(opcje.kanal);
  const stan = utworzStanDebaty(zrodlo, rejestr, { okno: opcje.okno });
  const strumien = utworzStrumienWypowiedzi(stan);
  const zrodloStrumienia = utworzZrodloStrumieniaDebaty(opcje.kanal, () => stan.okno());

  const rama = utworzRameOkna({
    tytul: 'Przebieg debaty',
    rola: 'monitor',
    kod: KOD_PANELU_DEBATY,
    przeznaczenie:
      'Tura bieżąca, skład i wypowiedzi wielu modeli obok rozmowy — wypowiedź rośnie na żywo.',
    modul: opcje.modul,
    przedrostek: opcje.przedrostek,
  });
  // Klasa modułu obok klas gospodarza: wygląd głosów niesie arkusz Roundtable
  // (`panel-debaty.css`), a nie arkusz modułu, w którego stosie panel stanął.
  rama.element.classList.add('dr-panel');

  const tresc = utworzStanTresci();
  let turaWidziana = stan.tura();

  function rysuj(): void {
    rysujPanelDebaty(stan, strumien, tresc);
    rama.ustawZnacznik(...znacznikGlosow(stan, strumien));
  }

  /**
   * Zmiana stanu debaty. Zmiana tury czyści gromadzenie strumienia: fragmenty
   * należą do wypowiedzi jednej tury, a `StanDebaty` też porzuca wtedy wykaz
   * wypowiedzi. Zostawione dokleiłyby zdania tury poprzedniej do mówców tury
   * nowej — czyli przypisałyby uczestnikowi słowa, których w niej nie powiedział.
   */
  function przyZmianieStanu(): void {
    if (stan.tura() !== turaWidziana) {
      turaWidziana = stan.tura();
      strumien.wyczysc();
    }
    rysuj();
  }

  function odswiez(): void {
    rejestr.odswiez();
    rysuj();
    tresc.potwierdzenie(
      powodBrakuObslugi(
        Command.RoundtableDebateGet,
        'Odświeżono wyłącznie wykaz kanałów. Przebieg debaty odczytuje ta komenda, ale panel jej jeszcze nie wywołuje, więc pokazuje to, co usłyszał od otwarcia.',
      ),
      true,
    );
  }

  const odswiezPrzycisk = przyciskAkcji('Odśwież nazwy kanałów', 'dn-btn dn-btn--atrament');
  odswiezPrzycisk.addEventListener('click', odswiez);
  rama.akcje.append(
    odswiezPrzycisk,
    przyciskBezKomendy(
      'Wczytaj przebieg tury',
      powodBrakuObslugi(
        Command.RoundtableDebateGet,
        'Skład, tury i wypowiedzi oddaje ta komenda jednym wywołaniem. Do czasu zbudowania jej obsługi ' +
          'panel pokazuje wyłącznie to, co zdarzenia przyniosły od jego otwarcia.',
      ),
    ),
  );
  rama.cialo.append(tresc.element);

  const odsubskrybujStan = stan.naZmiane(przyZmianieStanu);
  const odsubskrybujStrumien = zrodloStrumienia.naFragmentWypowiedzi((fragment) => {
    strumien.przyjmij(fragment);
    rysuj();
  });

  // Wykaz kanałów zamawia panel sam, bo bez nazw kanałów tożsamości uczestników
  // czytałyby się identyfikatorami. Pierwszy rysunek idzie od razu: bez niego
  // panel stałby pusty aż do pierwszego zdarzenia rdzenia, a stan pusty ma
  // własne zdanie mówiące, dlaczego jest pusty.
  rejestr.odswiez();
  rysuj();

  return {
    element: rama.element,
    odswiez,
    zamknij() {
      odsubskrybujStrumien();
      odsubskrybujStan();
      stan.zamknij();
    },
  };
}

/**
 * Plakietka nagłówka: ile głosów rośnie w tej chwili.
 *
 * Panel stoi w wąskiej kolumnie i bywa przewinięty, więc chwila „ktoś właśnie
 * mówi" musi być widoczna z samego nagłówka, bez zjeżdżania do wypowiedzi.
 * Zerowa liczba nie chowa plakietki na rzecz ciszy — mówi o niej wprost, dopóki
 * tura jest otwarta.
 */
function znacznikGlosow(
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
): [string, 'neutralna' | 'sukces' | 'ostrzezenie'] {
  let mowiacych = 0;
  let zerwanych = 0;
  for (const uczestnik of stan.uczestnicy()) {
    const glos = strumien.glos(uczestnik.id);
    if (glos === null) continue;
    if (glos.przyczyna !== '') zerwanych += 1;
    else if (!glos.domkniety) mowiacych += 1;
  }
  if (zerwanych > 0) return [`zerwane strumienie: ${zerwanych}`, 'ostrzezenie'];
  if (mowiacych > 0) return [`mówi teraz: ${mowiacych}`, 'sukces'];
  if (stan.definicjaTury() === null) return ['', 'neutralna'];
  return ['nikt nie mówi', 'neutralna'];
}
