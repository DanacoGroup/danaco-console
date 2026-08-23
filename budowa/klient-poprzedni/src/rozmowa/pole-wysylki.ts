import { elementIkony } from '../ikony/ikony';
import { NAPISY } from './etykiety-rozmowy';
import { opisStanu, type StanRozmowy } from './stan-rozmowy';
import type { WykazPoUkosniku } from './wykaz-po-ukosniku';

/** Zdarzenia pola wysyłki. */
export interface ObslugaWysylki {
  /** Operator zatwierdził wypowiedź. */
  naWyslanie(tekst: string): void;
  /** Operator zażądał przerwania tury. */
  naPrzerwanie(): void;
}

/** Ustawienia pola, których nie da się wyprowadzić z samej obsługi wysyłki. */
export interface OpcjePolaWysylki {
  /**
   * Gotowy pasek zlecenia stawiany w rzędzie pod polem wypowiedzi.
   *
   * Pole jest opcjonalne, bo stery paska żądają kompletu sterowania okna
   * (migawka, subskrypcja, `window.update`, `config.set`), a ten powstaje
   * dopiero w powłoce — stanowisko podglądu rozmowy (`podglad-rozmowy.ts`) go
   * nie ma i ma dalej działać. Pole przepuszcza element nietknięty: nie zna
   * ani nastaw, ani kontraktu.
   *
   * Pasek stoi pod polem, w tym samym wierszu co przycisk zatrzymania, bo
   * Operator ustawia na nim otoczenie zadania, zanim treść pójdzie do modelu.
   */
  pasekZlecenia?: HTMLElement;
  /**
   * Gotowy wykaz po ukośniku stawiany nad polem wypowiedzi.
   *
   * Pole go nie buduje i nie wie, skąd biorą się pozycje — zna wyłącznie cztery
   * ruchy klawiatury, które nim sterują (`otworz`, `ustawFraze`, `przesun`,
   * `zatwierdz`). Wykaz potrzebuje kanału i wątku rozmowy, a tych pole nie ma;
   * podaje go widok rozmowy, tak samo jak pasek zlecenia. Pominięty znaczy „to
   * stanowisko wykazu nie ma" — podgląd rozmowy (`podglad-rozmowy.ts`) nie ma
   * kanału i ma dalej działać.
   */
  wykazUkosnika?: WykazPoUkosniku;
}

/** Pole wpisywania wraz z przyciskami wysyłki i zatrzymania. */
export interface PoleWysylki {
  /** Element montowany w widoku rozmowy. */
  element: HTMLElement;
  /** Odświeża wskaźnik stanu wysyłania. */
  pokazStan(stan: StanRozmowy): void;
  /** Ustawia ognisko na polu wpisywania. */
  ustawOgnisko(): void;
  /** Dokłada gotowe polecenie do treści wypowiedzi; nie wysyła go. */
  wstaw(tekst: string): void;
  /**
   * Pokazuje w polu wypowiedź, którą wpisało do okna inne połączenie konta.
   *
   * Treść pojawia się od razu w całości, na krótko, i sama znika: wypowiedź
   * jest już wysłana — rdzeń rozgłosił ją zdarzeniem po przyjęciu komendy —
   * więc literowanie jej znak po znaku pokazywałoby zdanie niepełne.
   *
   * Odbitka nie zabiera pracy Operatorowi. Pole z ogniskiem albo z choćby
   * jednym znakiem treści nie jest ruszane, odbitka nie ma ogniska, nie da się
   * jej wysłać Enterem i nie wchodzi do treści — pole zostaje puste i gotowe
   * przez cały czas jej trwania.
   */
  odbijWypowiedzZZewnatrz(tekst: string): void;
}

/**
 * Jak długo odbitka cudzej wypowiedzi stoi w polu.
 *
 * Operator ma zdążyć zobaczyć, że tekst pojawił się w jego polu, i przeczytać
 * początek zdania, zanim zniknie. Krócej byłoby mrugnięciem nie do złapania
 * wzrokiem, dłużej — zawadą stojącą w miejscu pracy. Odbitka niczego nie
 * opóźnia: wypowiedź poszła do rdzenia, zanim to połączenie w ogóle się o niej
 * dowiedziało.
 */
const CZAS_ODBITKI_MS = 1200;

/**
 * Buduje pole wypowiedzi Operatora.
 *
 * Ani pole, ani przycisk wysyłki, ani przycisk zatrzymania nie tracą
 * klikalności w żadnym stanie. Przycisk zatrzymania jest czynny także wtedy,
 * gdy nic nie biegnie — rdzeń odpowiada wtedy `stopped` równym fałsz, a nie
 * błędem. Niegotowość opisuje wskaźnik obok przycisków, nie odebranie
 * możliwości działania.
 */
export function utworzPoleWysylki(
  obsluga: ObslugaWysylki,
  opcje: OpcjePolaWysylki = {},
): PoleWysylki {
  const element = document.createElement('form');
  // Wiersz dolny okna: wyściółka 8/16/12 px, tło panelu, bez kreski górnej.
  // Budowę niesie `.dn-rozmowa-dol` z biblioteki, `dc-wysylka` jest uchwytem
  // dla przejazdu po interfejsie.
  element.className = 'dn-rozmowa-dol dc-wysylka';

  // Powłoka pola — `.dn-prompt` z `komponenty/drobne.css`. Ramka, promień 8 px
  // i ognisko należą do powłoki, nie do obszaru tekstu: obrys ogniska obejmuje
  // wtedy grot, pole i przycisk jako jedno wejście.
  const powloka = document.createElement('div');
  powloka.className = 'dn-prompt';

  // Grot ❯ — sygnatura wejścia. Znak dekoracyjny, więc poza drzewem
  // dostępności; treść wejścia niesie etykieta pola.
  const grot = document.createElement('span');
  grot.className = 'dn-prompt-grot';
  grot.setAttribute('aria-hidden', 'true');
  grot.textContent = '❯';

  const pole = document.createElement('textarea');
  // Obszar tekstu bez własnej ramki i bez zmiany rozmiaru — rośnie od 40 px
  // do 160 px wewnątrz powłoki (`.dn-prompt-obszar`). `rows` jest jeden, bo
  // wysokość ustala `min-height`, nie liczba wierszy.
  pole.className = 'dn-prompt-obszar dc-wysylka__pole';
  pole.rows = 1;
  pole.placeholder = NAPISY.polePodpowiedz;
  pole.setAttribute('aria-label', NAPISY.polePodpowiedz);

  const pasek = document.createElement('div');
  // Pasek akcji pod polem: odstęp 4 px, odsunięcie od pola 8 px.
  pasek.className = 'dn-rozmowa-akcje dc-wysylka__pasek';

  const wskaznik = document.createElement('span');
  wskaznik.className = 'dc-wysylka__wskaznik';
  wskaznik.setAttribute('role', 'status');

  // Podpowiedź wykazu po ukośniku. Stoi w rzędzie akcji, a nie w `placeholder`:
  // po wpisaniu „/" pole nie jest puste, więc placeholder znika z ekranu
  // dokładnie w chwili, w której podpowiedź ma się pojawić. To ten sam wiersz
  // sterów co wskaźnik stanu i znika razem z zamknięciem wykazu.
  const podpowiedzWykazu = document.createElement('span');
  podpowiedzWykazu.className = 'dc-wysylka__podpowiedz';
  podpowiedzWykazu.setAttribute('role', 'status');
  podpowiedzWykazu.hidden = true;
  podpowiedzWykazu.textContent = NAPISY.wykazPodpowiedz;

  const przerwij = przycisk(NAPISY.przerwij, 'zatrzymaj', 'dn-btn dn-btn--zarys');
  const wyslij = przycisk(NAPISY.wyslij, 'wyslij', 'dn-btn dn-btn--atrament');
  wyslij.type = 'submit';

  // Odbitka wypowiedzi spoza tego połączenia. Leży nad polem jako warstwa
  // własna, poza treścią pola i poza drzewem wejścia — dzięki temu pokazanie
  // cudzego zdania nie może zamienić się w wysłanie go przez Operatora.
  // `aria-live` niesie ją czytnikowi ekranu.
  const odbitka = document.createElement('output');
  odbitka.className = 'dc-wysylka__odbitka';
  odbitka.hidden = true;
  odbitka.setAttribute('aria-live', 'polite');
  let zegarOdbitki = 0;

  powloka.append(grot, pole, odbitka, wyslij);
  // Pasek zlecenia idzie pierwszy w rzędzie. Wskaźnik stanu ma
  // `margin-right: auto`, więc odsuwa przycisk zatrzymania do prawej krawędzi;
  // stery wstawione za nim wylądowałyby po tamtej stronie rzędu, razem
  // z zatrzymaniem tury.
  if (opcje.pasekZlecenia !== undefined) pasek.append(opcje.pasekZlecenia);
  pasek.append(podpowiedzWykazu, wskaznik, przerwij);
  // Wykaz po ukośniku idzie przed powłoką pola, bo otwiera się w górę: pole
  // zostaje na miejscu, a wykaz zajmuje wolną przestrzeń rozmowy nad nim.
  // Element jest pływający i w rzędzie nie zajmuje nic.
  const wykaz = opcje.wykazUkosnika;
  if (wykaz !== undefined) element.append(wykaz.element);
  element.append(powloka, pasek);

  function zatwierdz(): void {
    const tresc = pole.value;
    if (tresc.trim().length === 0) return;
    pole.value = '';
    obsluga.naWyslanie(tresc);
  }

  /* Wykaz po ukośniku — wyzwalacz z klawiatury.

     Wykaz pojawia się natychmiast po wpisaniu „/", nie po dopisaniu liter.
     Wyzwalaczem jest zdarzenie `input`, a nie `keydown`: znak jest już w polu,
     więc warunek „ukośnik stoi na początku treści" da się sprawdzić na treści,
     a nie zgadywać z klawisza (martwe klawisze, wklejenie, IME).

     Ukośnik w środku zdania nie otwiera niczego — ścieżka `/home/ubuntu`
     wpisana w wypowiedź jest tekstem, nie komendą. */

  /** Czy treść pola jest wpisem po ukośniku (i jaka jest jego fraza). */
  function frazaUkosnika(tresc: string): string | null {
    if (!tresc.startsWith('/')) return null;
    const reszta = tresc.slice(1);
    // Nowy wiersz kończy komendę: to już wypowiedź wielowierszowa, nie wpis.
    return reszta.includes('\n') ? null : reszta;
  }

  function odswiezPodpowiedz(): void {
    podpowiedzWykazu.hidden = wykaz === undefined || !wykaz.otwarty();
  }

  if (wykaz !== undefined) {
    wykaz.naZmiane(() => odswiezPodpowiedz());
    pole.addEventListener('input', () => {
      const fraza = frazaUkosnika(pole.value);
      if (fraza === null) {
        // Skasowanie ukośnika zwija wykaz — tak samo jak Escape.
        if (wykaz.otwarty()) wykaz.zamknij();
        odswiezPodpowiedz();
        return;
      }
      if (!wykaz.otwarty()) wykaz.otworz();
      wykaz.ustawFraze(fraza);
      odswiezPodpowiedz();
    });
  }

  element.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    zatwierdz();
  });

  przerwij.addEventListener('click', () => obsluga.naPrzerwanie());

  // Enter wysyła, Shift+Enter dokłada wiersz — zwyczaj konsoli, nie edytora.
  //
  // Przy otwartym wykazie te same klawisze znaczą co innego: strzałki wędrują
  // po pozycjach zamiast po tekście, Enter zatwierdza wyróżnioną pozycję
  // zamiast wysyłać wypowiedź, Escape zwija wykaz. Zamiana obowiązuje wyłącznie
  // wtedy, gdy wykaz stoi na ekranie.
  pole.addEventListener('keydown', (zdarzenie) => {
    if (wykaz !== undefined && wykaz.otwarty()) {
      if (zdarzenie.key === 'ArrowDown' || zdarzenie.key === 'ArrowUp') {
        zdarzenie.preventDefault();
        wykaz.przesun(zdarzenie.key === 'ArrowDown' ? 1 : -1);
        return;
      }
      if (zdarzenie.key === 'Escape') {
        zdarzenie.preventDefault();
        zdarzenie.stopPropagation();
        wykaz.zamknij();
        odswiezPodpowiedz();
        return;
      }
      if (zdarzenie.key === 'Enter' && !zdarzenie.shiftKey) {
        zdarzenie.preventDefault();
        // Pole czyści się po zatwierdzeniu, bo wybór pozycji nie jest
        // wypowiedzią: nic nie idzie do modelu, więc nie ma czego zostawiać
        // w polu ani czego wysyłać.
        const wybrano = wykaz.zatwierdz();
        odswiezPodpowiedz();
        if (wybrano) {
          pole.value = '';
          return;
        }
        // Wykaz przycięty frazą do zera nie ma czego zatwierdzić, więc treść
        // pola idzie zwykłą drogą wysyłki zamiast zniknąć bez śladu.
        zatwierdz();
        return;
      }
    }
    if (zdarzenie.key !== 'Enter' || zdarzenie.shiftKey) return;
    zdarzenie.preventDefault();
    zatwierdz();
  });

  return {
    element,
    pokazStan(stan) {
      wskaznik.textContent = opisStanu(stan);
      // `aria-busy` niesie stan pracy bez odbierania klikalności.
      wyslij.setAttribute('aria-busy', String(stan.wysyla));
      element.dataset['wysyla'] = String(stan.wysyla);
    },
    ustawOgnisko: () => pole.focus(),

    odbijWypowiedzZZewnatrz(tekst) {
      const tresc = tekst.trim();
      if (tresc.length === 0) return;
      // Operator pisze — pole należy do niego i nikt go z niego nie wyrywa.
      if (pole.value.length > 0 || document.activeElement === pole) return;

      // Odbitka mieszka w warstwie nad polem, nie w jego treści: wpisana do
      // `pole.value` poszłaby drugi raz do rdzenia, gdyby Operator nacisnął
      // Enter w tej samej sekundzie.
      odbitka.textContent = tresc;
      odbitka.hidden = false;
      powloka.dataset['odbitka'] = 'tak';
      window.clearTimeout(zegarOdbitki);
      zegarOdbitki = window.setTimeout(() => {
        odbitka.hidden = true;
        odbitka.textContent = '';
        delete powloka.dataset['odbitka'];
      }, CZAS_ODBITKI_MS);
    },

    // Narzędzie paska promptu i akcja modułu wstawiają polecenie, a nie
    // wysyłają je: decyzja o turze zostaje przy Operatorze, który może
    // polecenie uzupełnić — część pozycji kończy się dwukropkiem właśnie po to.
    // Wstawienie dokłada się do treści zastanej, zamiast ją kasować.
    wstaw(tekst) {
      const zastana = pole.value;
      const rozdzielnik = zastana.length > 0 && !zastana.endsWith('\n') ? '\n' : '';
      pole.value = `${zastana}${rozdzielnik}${tekst}`;
      pole.focus();
      pole.setSelectionRange(pole.value.length, pole.value.length);
    },
  };
}

/** Przycisk paska wysyłki: ikona i etykieta słowna. */
function przycisk(etykieta: string, ikona: 'wyslij' | 'zatrzymaj', klasa: string): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.append(elementIkony(ikona, { rozmiar: 16, klasa: 'dn-ikona' }));

  const napis = document.createElement('span');
  napis.textContent = etykieta;
  element.append(napis);
  return element;
}
