import type { PracaAsystenta } from '../aplikacja/zrodlo-posuniec';
import { elementIkony } from '../ikony/ikony';

/** Pływający favikon Asystenta — kontrolka obecna w każdym module, byt powłoki, a nie samego modułu roboczego. */
export interface FavikonPlywajacy {
  element: HTMLElement;
  /** Odzwierciedla, czy dymek jest odsłonięty — dla `aria-expanded`. */
  ustawOtwarty(otwarty: boolean): void;
  // Nanosi pracę asystenta w toku albo jej brak; `null` jest odpowiedzią, nie brakiem odpowiedzi.
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

/** Zdanie o pracy — jedno źródło dla `aria-label`, `title`, znaku stanu i wiersza w dymku, złożone znak w znak. */
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

  // Kropka stoi bez względu na stan — kształt niesie znaczenie obok barwy, nigdy sama barwa.
  const kropka = document.createElement('span');
  kropka.className = 'ap-favikon__kropka';
  kropka.dataset.glos = String(opcje.glosDziala);
  if (opcje.glosDziala) kropka.append(elementIkony('mikrofon', { rozmiar: 14 }));
  else kropka.textContent = '×';
  kropka.setAttribute('aria-hidden', 'true');

  // Znak pracy stoi tylko wtedy, gdy asystent naprawdę prowadzi zlecenie.
  const praca = document.createElement('span');
  praca.className = 'ap-favikon__praca';
  praca.hidden = true;
  praca.textContent = '▶';
  praca.setAttribute('aria-hidden', 'true');

  // Licznik nieprzeczytanych posunięć jest liczbą, nie kropką — inna wiadomość dla Operatora.
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
        // Powyżej dziewięciu liczba przestaje cokolwiek zmieniać, a rozpycha przycisk.
      licznik.textContent = nieprzeczytane > 9 ? '9+' : String(nieprzeczytane);
      licznik.hidden = nieprzeczytane === 0;
      element.dataset.nieprzeczytane = String(nieprzeczytane);
      opisz();
    },
  };
}

/** Odmiana rzeczownika po liczbie posunięć — etykieta czytana na głos ma brzmieć naturalnie po polsku dla Operatora. */
function slowoPosuniecia(ile: number): string {
  if (ile === 1) return 'posunięcie asystenta';
  const reszta = ile % 10;
  const setka = ile % 100;
  if (reszta >= 2 && reszta <= 4 && (setka < 12 || setka > 14)) return 'posunięcia asystenta';
  return 'posunięć asystenta';
}
