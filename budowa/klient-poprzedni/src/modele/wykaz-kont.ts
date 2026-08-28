import type { Account } from '../../../shared/contract';
import { nazwaKrotka } from './rodzaje-kont';
import type { StanKont } from './stan-kont';
import { znakWykazu } from './znak-wykazu';

/**
 * Wykaz kont tworzy lewą kolumnę sekcji kont. Każdy wiersz podaje nazwę konta, jego
 * rodzaj, dostawcę oraz stan wyrażony plakietkami, a pusty rejestr zastępuje wiersze
 * zdaniem dobranym do fazy odczytu.
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

/**
 * Buduje jeden wiersz wykazu: nazwę konta, plakietki stanu oraz drugą linię z dostawcą
 * i modelem domyślnym. Kliknięcie wiersza czyni konto wyborem czynnym stanu sekcji.
 */
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

/**
 * Składa drugą linię wiersza z krótkiej nazwy rodzaju konta, nazwy dostawcy oraz modelu
 * domyślnego, gdy konto go wskazuje; człony puste zostają pominięte.
 */
function podpis(konto: Account): string {
  const czlony = [nazwaKrotka(konto.kind), konto.provider];
  if (konto.defaultModel !== undefined && konto.defaultModel !== '') {
    czlony.push(konto.defaultModel);
  }
  return czlony.filter((czlon) => czlon !== '').join(' · ');
}

/**
 * Tworzy pojedynczy znak wiersza wykazu kont, opierając go na plakietce biblioteki
 * komponentów osadzonej w miejscu przeznaczonym na znaki wiersza.
 */
function znak(tresc: string, klasa: string): HTMLElement {
  return znakWykazu(tresc, klasa, 'dm-wykaz__znak');
}

/**
 * Składa plakietki stanu konta: oznaczenie konta domyślnego, obecność albo brak
 * zapisanego poświadczenia oraz oznaczenie konta nieczynnego.
 */
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
 * Buduje stan pusty wykazu na klasie pustego stanu z biblioteki komponentów, dobierając
 * zdanie opisu do fazy odczytu rejestru kont.
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
