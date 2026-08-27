import { przyciskBezKomendy } from '../modele/kontrolki-formularza-braki';
import { utworzStanTresci } from '../komponenty/stan-tresci';
import type { ObszarIzolacji, ZaleznosciObszaru } from './obszary';

/** Obszar Oś rozstrzygania — trzy osie zapisu punktu izolacji w kolejności: konto, model, platforma domyślna. */
export function utworzObszar(_zaleznosci: ZaleznosciObszaru): ObszarIzolacji {
  const tresc = utworzStanTresci('pi');

  function odswiez(): void {
    tresc.tresc().append(zbudujWidok());
  }

  odswiez();

  return { element: tresc.element, odswiez };
}

/** Jedna z trzech osi rozstrzygania: nazwa po polsku oraz zdanie, co znaczy zapis wartości na danej osi. */
interface DefinicjaOsi {
  nazwa: string;
  zdanie: string;
}

/**
 * Trzy osie w kolejności rozstrzygania w ramach poziomu — dosłownie
 * `osieOdNajwezszej` z `osie.go`: konto (byt konkretny) → model (klasa) →
 * platforma (tło, wartość domyślna gdy oś pominięta).
 */
const OSIE_W_KOLEJNOSCI: readonly DefinicjaOsi[] = [
  {
    nazwa: 'Konto',
    zdanie:
      'Zapis obowiązuje wyłącznie wskazane konto Operatora. W obrębie tego samego poziomu zasięgu ' +
      'to byt najbardziej celowy tej osi — przeważa nad zapisem na modelu i na platformie.',
  },
  {
    nazwa: 'Model',
    zdanie:
      'Zapis obowiązuje całą klasę modelu (np. rodzinę modelu kanału), niezależnie od tego, które ' +
      'konto go używa. W obrębie poziomu ustępuje zapisowi na koncie, ale przeważa nad tłem platformy.',
  },
  {
    nazwa: 'Platforma',
    zdanie:
      'Zapis obowiązuje jako tło — całą platformę, gdy w obrębie tego samego poziomu nic węższego ' +
      '(konto, model) nie zostało wskazane. Oś pominięta domyślnie znaczy platformę — pusty wybór nigdy nie jest błędem.',
  },
];

function zbudujWidok(): HTMLElement {
  const element = document.createElement('div');
  element.className = 'pi-lista';

  element.append(zbudujWprowadzenie(), zbudujRelacjeDoPoziomu(), zbudujListeOsi(), zbudujMiejsceZapisu());
  return element;
}

function zbudujWprowadzenie(): HTMLElement {
  const wstep = document.createElement('p');
  wstep.className = 'pi-osie__wstep';
  wstep.textContent =
    'Każdy z jedenastu punktów izolacji zapisuje się pod dwiema współrzędnymi: poziomem zasięgu ' +
    '(JAK WĄSKO) i osią rozstrzygania (DLA CZEGO). Ten obszar pokazuje drugą współrzędną — oś — oraz ' +
    'to, że rozstrzyga się ona dopiero wewnątrz poziomu już wybranego w obszarze „Poziom zasięgu".';
  return wstep;
}

/** Karta pokazująca, że poziom i oś to dwie prostopadłe współrzędne izolacji, a nie jedna wspólna drabina. */
function zbudujRelacjeDoPoziomu(): HTMLElement {
  const karta = document.createElement('div');
  karta.className = 'dn-karta';

  const naglowek = document.createElement('div');
  naglowek.className = 'dn-karta-naglowek';
  const tytul = document.createElement('span');
  tytul.className = 'dn-karta-tytul';
  tytul.textContent = 'Dwie współrzędne, nie jedna drabina';
  naglowek.append(tytul);

  const cialo = document.createElement('div');
  cialo.className = 'dn-karta-cialo pi-karta__cialo';

  const wiersze = document.createElement('div');
  wiersze.className = 'pi-osie__wiersze';
  wiersze.append(
    zbudujPlakietke('1 · rozstrzyga PIERWSZY', false),
    zbudujOpisWspolrzednej(
      'Poziom zasięgu (obszar „Poziom zasięgu") — JAK WĄSKO: okno → rola → sesja → projekt → ' +
        'para modułów → moduł → środowisko → globalny.',
    ),
    zbudujPlakietke('2 · rozstrzyga W RAMACH poziomu', true),
    zbudujOpisWspolrzednej('Oś rozstrzygania (ten obszar) — DLA CZEGO: konto → model → platforma.'),
  );

  const wyjasnienie = document.createElement('p');
  wyjasnienie.className = 'pi-osie__wyjasnienie';
  wyjasnienie.textContent =
    'Poziom rozstrzyga do końca, zanim oś w ogóle wejdzie w grę: zapis na koncie na poziomie globalnym ' +
    'NIE bije zapisu zrobionego wprost w oknie komunikacji. Dopiero wewnątrz zwycięskiego poziomu oś ' +
    'wskazuje adresata — dla kogo ta wartość obowiązuje. Wybór osi nigdy nie unieważnia wyboru poziomu.';

  cialo.append(wiersze, wyjasnienie);
  karta.append(naglowek, cialo);
  return karta;
}

function zbudujPlakietke(etykieta: string, sygnal: boolean): HTMLElement {
  const plakietka = document.createElement('span');
  plakietka.className = sygnal ? 'dn-plakietka dn-plakietka--sygnal' : 'dn-plakietka';
  plakietka.textContent = etykieta;
  return plakietka;
}

function zbudujOpisWspolrzednej(zdanie: string): HTMLElement {
  const opis = document.createElement('span');
  opis.className = 'pi-osie__opis-wspolrzednej';
  opis.textContent = zdanie;
  return opis;
}

/** Trzy osie w kolejności rozstrzygania — lista uporządkowana, bo ta kolejność jest treścią, nie ozdobą. */
function zbudujListeOsi(): HTMLElement {
  const lista = document.createElement('ol');
  lista.setAttribute('aria-label', 'Trzy osie rozstrzygania w kolejności, w jakiej rozstrzygają się w ramach poziomu');
  lista.className = 'pi-osie__lista';

  OSIE_W_KOLEJNOSCI.forEach((os, indeks) => {
    const pozycja = document.createElement('li');
    pozycja.append(zbudujKarteOsi(os, indeks + 1));
    lista.append(pozycja);
  });

  return lista;
}

function zbudujKarteOsi(os: DefinicjaOsi, kolejnosc: number): HTMLElement {
  const karta = document.createElement('div');
  karta.className = 'dn-karta';

  const naglowek = document.createElement('div');
  naglowek.className = 'dn-karta-naglowek';

  const plakietka = document.createElement('span');
  plakietka.className = 'dn-plakietka';
  plakietka.title = `Miejsce ${String(kolejnosc)} z 3 w kolejności rozstrzygania osi w ramach poziomu`;
  plakietka.textContent = String(kolejnosc);

  const tytul = document.createElement('span');
  tytul.className = 'dn-karta-tytul';
  tytul.textContent = os.nazwa;

  naglowek.append(plakietka, tytul);

  const cialo = document.createElement('div');
  cialo.className = 'dn-karta-cialo';
  cialo.textContent = os.zdanie;

  karta.append(naglowek, cialo);
  return karta;
}

/** Jedyna czynność tego obszaru: wybór bytu wskazanej osi i zapis wartości na tej osi rozstrzygania izolacji. */
function zbudujMiejsceZapisu(): HTMLElement {
  const karta = document.createElement('div');
  karta.className = 'dn-karta';

  const naglowek = document.createElement('div');
  naglowek.className = 'dn-karta-naglowek';
  const tytul = document.createElement('span');
  tytul.className = 'dn-karta-tytul';
  tytul.textContent = 'Zapis wartości na wskazanej osi';
  naglowek.append(tytul);

  const cialo = document.createElement('div');
  cialo.className = 'dn-karta-cialo pi-karta__cialo';

  const opis = document.createElement('p');
  opis.className = 'pi-osie__opis';
  opis.textContent =
    'Docelowo: wybór bytu tej osi (identyfikator konta albo modelu) i zapis wartości dla klucza ' +
    'wskazanego w obszarze „Kontekst" lub „Zakres techniczny", na poziomie wskazanym w obszarze ' +
    '„Poziom zasięgu".';

  const powod =
    'Kontrakt nie ma komendy zapisu na OSI rozstrzygania. Dwanaście komend isolation.* jest wpiętych ' +
    'i okno je woła, ale żądania zapisu (isolation.context.set, isolation.technical.set) niosą wyłącznie ' +
    'poziom zasięgu, byt poziomu i warstwę — pola osi (konto/model/platforma) w nich nie ma. ' +
    'Brakuje komendy w kontrakcie, nie podłączenia w tym oknie.';
  const przycisk = przyciskBezKomendy('Podgląd i zapis na tej osi', powod);

  cialo.append(opis, przycisk);
  karta.append(naglowek, cialo);
  return karta;
}
