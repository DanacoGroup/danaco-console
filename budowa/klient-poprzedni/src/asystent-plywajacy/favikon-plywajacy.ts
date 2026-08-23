import type { PracaAsystenta } from '../aplikacja/zrodlo-posuniec';
import { elementIkony } from '../ikony/ikony';

/**
 * Pływający favikon Asystenta — kontrolka obecna w każdym module.
 *
 * Favikon jest bytem powłoki, nie modułu: nie wchodzi na scenę okien
 * równoległych, nie liczy się do sufitu `LICZBA_MAX` i nie znika przy zmianie
 * modułu. Jedna odpowiedzialność: przycisk i jego stan — rozmowy i drogi do
 * rdzenia leżą w `okno-dymkowe.ts` i `stan-dymka.ts`.
 *
 * Favikon jest sterem, nie wyświetlaczem, więc etykieta dostępności niesie
 * wartość bieżącą („Asystent pracuje: <tytuł zlecenia> · etap 2 z 5"), a nie
 * napis rodzajowy „Asystent".
 *
 * Licznik nieprzeczytanych istnieje, bo dymek bywa zwinięty: jest przywoływany,
 * a nie rysowany z urzędu, więc asystent potrafi wykonać ciąg posunięć, zanim
 * Operator go otworzy. Bez licznika ciąg ten nie zostawiałby na ekranie śladu.
 *
 * Kropka stanu mówi o głosie, nie o łączności — wskaźnik łączności z rdzeniem
 * stoi w pasku górnym (`aplikacja/wskaznik-lacznosci.ts`) i drugiego się nie
 * stawia. Ta kropka niesie to, czego nie niesie nic innego: czy kanał głosowy
 * modułu w ogóle istnieje. Kanał zbudowany oznacza ikona zestawu; odpowiedź
 * przeciwna zostaje znakiem typograficznym, bo zestaw nie niesie mikrofonu
 * przekreślonego, a własnego znaku dorysować nie wolno. Znak pracy stoi obok
 * kropki głosu i jest od niej niezależny.
 */

export interface FavikonPlywajacy {
  element: HTMLElement;
  /** Odzwierciedla, czy dymek jest odsłonięty — dla `aria-expanded`. */
  ustawOtwarty(otwarty: boolean): void;
  /**
   * Nanosi pracę asystenta w toku albo jej brak.
   *
   * `null` znaczy „asystent nie prowadzi zlecenia" i tak też jest napisane —
   * pusty stan nie jest tu brakiem odpowiedzi, tylko odpowiedzią.
   */
  ustawPrace(praca: PracaAsystenta | null): void;
  /** Ile posunięć asystenta czeka nieprzeczytanych w zwiniętym dymku. */
  ustawNieprzeczytane(ile: number): void;
}

export interface OpcjeFavikonu {
  /** Czy kanał głosowy jest zbudowany; rozstrzyga kropkę i podpowiedź. */
  glosDziala: boolean;
  /** Naciśnięcie favikonu — przełączenie dymka. */
  naNacisniecie(): void;
}

/**
 * Zdanie o pracy — jedno źródło dla `aria-label`, `title`, znaku stanu
 * i wiersza w dymku.
 *
 * Składane znak w znak tak samo jak na pasie dolnym (`aplikacja/pas-posuniec.ts`,
 * funkcja `opiszPrace`), łącznie ze środkową kropką przed etapem. Ta sama praca
 * pokazana w trzech miejscach ma się czytać jednakowo — różnica choćby
 * w przecinku każe Operatorowi sprawdzać, czy to na pewno to samo zlecenie.
 */
export function zdanieOPracy(praca: PracaAsystenta | null): string {
  if (praca === null) return 'Asystent nie prowadzi zlecenia';
  const etapy = praca.etapow > 0 ? ` · etap ${praca.etap} z ${praca.etapow}` : '';
  return `Asystent pracuje: ${praca.tytul}${etapy}`;
}

export function utworzFavikonPlywajacy(opcje: OpcjeFavikonu): FavikonPlywajacy {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'ap-favikon';
  element.dataset.glos = String(opcje.glosDziala);
  element.dataset.pracuje = 'nie';
  element.setAttribute('aria-expanded', 'false');

  const znak = document.createElement('span');
  znak.className = 'ap-favikon__znak';
  znak.append(elementIkony('rozmowa', { rozmiar: 24 }));

  // Kropka stoi bez względu na stan: „głos jest" i „głosu nie ma" to dwie
  // odpowiedzi, nie odpowiedź i jej brak. Kształt niesie znaczenie obok barwy,
  // bo stan nigdy nie opiera się na samej barwie.
  const kropka = document.createElement('span');
  kropka.className = 'ap-favikon__kropka';
  kropka.dataset.glos = String(opcje.glosDziala);
  if (opcje.glosDziala) kropka.append(elementIkony('mikrofon', { rozmiar: 14 }));
  else kropka.textContent = '×';
  kropka.setAttribute('aria-hidden', 'true');

  // Znak pracy — kształt, nie sama barwa. Stoi tylko wtedy, gdy asystent
  // naprawdę prowadzi zlecenie; pusty znak byłby drugą kropką bez znaczenia.
  const praca = document.createElement('span');
  praca.className = 'ap-favikon__praca';
  praca.hidden = true;
  praca.textContent = '▶';
  praca.setAttribute('aria-hidden', 'true');

  // Licznik nieprzeczytanych posunięć — liczba, nie kropka. „Coś się stało"
  // i „stało się sześć rzeczy" to dla Operatora dwie różne wiadomości.
  const licznik = document.createElement('span');
  licznik.className = 'ap-favikon__licznik';
  licznik.hidden = true;
  licznik.setAttribute('aria-hidden', 'true');

  element.append(znak, kropka, praca, licznik);
  element.addEventListener('click', opcje.naNacisniecie);

  let biezacaPraca: PracaAsystenta | null = null;
  let nieprzeczytane = 0;

  /** Składa etykietę z trzech prawd naraz: głos, praca, zaległości. */
  function opisz(): void {
    const oGlosie = opcje.glosDziala
      ? 'Asystent — rozmowa głosowa i tekstowa'
      : 'Asystent — rozmowa TEKSTOWA. Rdzeń nie niesie jeszcze mowy.';
    const oPracy = zdanieOPracy(biezacaPraca);
    const oZaleglosciach =
      nieprzeczytane > 0
        ? `${nieprzeczytane} ${slowoPosuniecia(nieprzeczytane)} do przeczytania w dymku`
        : '';

    const czesci = [oPracy, oZaleglosciach, oGlosie].filter((zdanie) => zdanie !== '');
    const pelne = czesci.join('. ');
    // Etykieta niesie wartość bieżącą, a dopiero po niej czynność przycisku.
    element.setAttribute(
      'aria-label',
      `${pelne}. ${element.dataset.otwarty === 'true' ? 'Zwiń' : 'Otwórz'} dymek rozmowy z Asystentem.`,
    );
    element.title = pelne;
  }

  opisz();

  return {
    element,

    ustawOtwarty(otwarty) {
      element.setAttribute('aria-expanded', String(otwarty));
      element.dataset.otwarty = String(otwarty);
      opisz();
    },

    ustawPrace(nowa) {
      biezacaPraca = nowa;
      element.dataset.pracuje = nowa === null ? 'nie' : 'tak';
      praca.hidden = nowa === null;
      opisz();
    },

    ustawNieprzeczytane(ile) {
      nieprzeczytane = Math.max(0, Math.trunc(ile));
      // Powyżej dziewięciu dokładna liczba przestaje cokolwiek zmieniać dla
      // Operatora, a rozpycha przycisk. „9+" jest tu odpowiedzią pełną.
      licznik.textContent = nieprzeczytane > 9 ? '9+' : String(nieprzeczytane);
      licznik.hidden = nieprzeczytane === 0;
      element.dataset.nieprzeczytane = String(nieprzeczytane);
      opisz();
    },
  };
}

/** Odmiana rzeczownika — etykieta czytana na głos ma brzmieć po polsku. */
function slowoPosuniecia(ile: number): string {
  if (ile === 1) return 'posunięcie asystenta';
  const reszta = ile % 10;
  const setka = ile % 100;
  if (reszta >= 2 && reszta <= 4 && (setka < 12 || setka > 14)) return 'posunięcia asystenta';
  return 'posunięć asystenta';
}
