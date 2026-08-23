import type { Account } from '../../../shared/contract';
import { nazwaKrotka } from './rodzaje-kont';
import type { StanKont } from './stan-kont';
import { znakWykazu } from './znak-wykazu';

/**
 * Wykaz kont — lewa kolumna sekcji kont.
 *
 * Wiersz mówi cztery rzeczy naraz: nazwę, rodzaj, dostawcę oraz stan konta
 * wyrażony plakietkami. Plakietka poświadczenia jest jedyną informacją o nim,
 * jaką klient ma prawo pokazać — kontrakt nie zwraca treści poświadczenia
 * żadną komendą.
 *
 * Wiersz konta nieczynnego nie jest wygaszony ani nieklikalny. Nieczynność
 * jest stanem danych, nie blokadą interfejsu — w takie konto trzeba móc wejść,
 * żeby je z powrotem uruchomić.
 *
 * Rejestr pusty nie daje pustego prostokąta — mówi wprost, że rdzeń nie oddał
 * ani jednego konta, i zostawia widok czynny.
 */
export interface WykazKont {
  /** Kolumna osadzana w panelu kont. */
  element: HTMLElement;
  /** Przebudowuje wykaz ze stanu i oznacza konto czynne. */
  odswiez(): void;
}

export function utworzWykazKont(stan: StanKont): WykazKont {
  const lista = document.createElement('div');
  lista.className = 'dm-wykaz__lista';

  const element = document.createElement('nav');
  element.className = 'dn-boczna dm-wykaz';
  element.setAttribute('aria-label', 'Wykaz kont modeli i kont code CLI');
  element.append(lista);

  function odswiez(): void {
    const konta = stan.konta();
    const czynne = stan.wybrane();

    if (konta.length === 0) {
      lista.replaceChildren(stanPusty(stan));
      return;
    }

    lista.replaceChildren(
      ...konta.map((konto) => wiersz(konto, konto.id === czynne?.id, stan)),
    );
  }

  return { element, odswiez };
}

/** Jeden wiersz wykazu: nazwa, plakietki stanu, dostawca i model domyślny. */
function wiersz(konto: Account, czynne: boolean, stan: StanKont): HTMLElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-boczna-pozycja dm-wykaz__pozycja';
  przycisk.dataset.konto = konto.id;
  if (czynne) przycisk.setAttribute('aria-current', 'true');
  przycisk.addEventListener('click', () => stan.wybierz(konto.id));

  const nazwa = document.createElement('span');
  nazwa.className = 'dm-wykaz__nazwa';
  nazwa.textContent = konto.name;

  const opis = document.createElement('span');
  opis.className = 'dm-wykaz__opis';
  opis.textContent = podpis(konto);

  const znaki = document.createElement('span');
  znaki.className = 'dm-wykaz__znaki';
  znaki.append(...plakietki(konto));

  przycisk.append(nazwa, opis, znaki);
  return przycisk;
}

/** Druga linia wiersza: dostawca oraz model domyślny, gdy konto go wskazuje. */
function podpis(konto: Account): string {
  const czlony = [nazwaKrotka(konto.kind), konto.provider];
  if (konto.defaultModel !== undefined && konto.defaultModel !== '') {
    czlony.push(konto.defaultModel);
  }
  return czlony.filter((czlon) => czlon !== '').join(' · ');
}

/** Znak wiersza wykazu kont — plakietka biblioteki w miejscu znaków wiersza. */
function znak(tresc: string, klasa: string): HTMLElement {
  return znakWykazu(tresc, klasa, 'dm-wykaz__znak');
}

/** Plakietki stanu konta: domyślność, poświadczenie, czynność. */
function plakietki(konto: Account): HTMLElement[] {
  const znaki: HTMLElement[] = [];

  if (konto.isDefault) {
    znaki.push(znak('domyślne', 'dn-plakietka dn-plakietka--sygnal'));
  }
  znaki.push(
    konto.hasCredential
      ? znak('poświadczenie zapisane', 'dn-plakietka dn-plakietka--sukces')
      : znak('bez poświadczenia', 'dn-plakietka dn-plakietka--ostrzezenie'),
  );
  if (!konto.enabled) {
    znaki.push(znak('nieczynne', 'dn-plakietka'));
  }

  return znaki;
}

/**
 * Stan pusty wykazu, oparty na klasie `.dn-pusty-stan` biblioteki komponentów.
 *
 * Zdanie dobiera się do fazy odczytu: ten sam pusty wykaz znaczy co innego
 * w trakcie pytania rdzenia, co innego po jego odmowie, a co innego, gdy rdzeń
 * odpowiedział i rejestr naprawdę nie ma konta. Jedno zdanie na trzy przypadki
 * nie rozstrzygałoby, czy czekać, czy działać.
 */
function stanPusty(stan: StanKont): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dn-pusty-stan dm-wykaz__pusto';

  const tytul = document.createElement('p');
  tytul.className = 'dn-pusty-stan-tytul';
  tytul.textContent = 'Brak kont modeli';

  const opis = document.createElement('p');
  opis.className = 'dn-pusty-stan-opis';
  opis.textContent = zdaniePustego(stan);

  element.append(tytul, opis);
  return element;
}

function zdaniePustego(stan: StanKont): string {
  switch (stan.faza()) {
    case 'spoczynek':
    case 'odczyt':
      return 'Rejestr kont jedzie z rdzenia. Wykaz pozostaje czynny.';
    case 'blad':
      return 'Rejestr kont nie dotarł z rdzenia — powód nad zakładkami sekcji. Odczyt można powtórzyć.';
    case 'gotowe':
      return 'Rdzeń odpowiedział, ale nie zna konta spełniającego warunek wykazu. Załóż konto formularzem obok albo zdejmij ograniczenie rodzaju.';
  }
}
