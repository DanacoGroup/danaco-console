// Panel wciąga swoje arkusze sam, bo gospodarzem bywa moduł, którego arkusz o module Roundtable nic nie wie.
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

/** Panel debaty — przebieg debaty obok rozmowy jako zwykły panel stosu; kod pozycji to ten sam, którym rejestr i wytwórnia paneli go zawołają. */
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
  // Klasa modułu obok klas gospodarza: wygląd głosów niesie arkusz Roundtable, nie moduł gospodarza.
  rama.element.classList.add('dr-panel');

  const tresc = utworzStanTresci();
  let turaWidziana = stan.tura();

  function rysuj(): void {
    rysujPanelDebaty(stan, strumien, tresc);
    rama.ustawZnacznik(...znacznikGlosow(stan, strumien));
  }

  // Zmiana tury czyści gromadzenie strumienia, żeby nie dokleić zdań tury poprzedniej do mówców nowej.
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

  // Wykaz kanałów zamawia panel sam, bo bez nazw tożsamości uczestników czytałyby się identyfikatorami.
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
 * Plakietka nagłówka: ile głosów rośnie w tej chwili, widoczna z samego nagłówka bez
 * zjeżdżania do wypowiedzi.
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
