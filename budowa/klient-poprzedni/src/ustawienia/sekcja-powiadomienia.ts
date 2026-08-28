import type { Kanal } from '../protokol/kanal';
import type { SekcjaUstawien } from './sekcje';
import {
  KLUCZ_GLOWNY,
  kluczKlasy,
  utworzZrodloPowiadomien,
  type KlasaZdarzenia,
  type StanKlasy,
  type StanPowiadomien,
  type ZrodloPowiadomien,
} from './zrodlo-powiadomien';

/** Sekcja powiadomień pokazuje przełącznik główny, siedem klas zdarzeń w tabeli i kanały dostarczenia, zapisując każdą zmianę do rdzenia natychmiast, bez przycisku zbiorczego. */
export function utworzSekcjePowiadomienia(kanal: Kanal): SekcjaUstawien {
  const zrodlo: ZrodloPowiadomien = utworzZrodloPowiadomien(kanal);

  const element = document.createElement('div');
  element.className = 'du-sekcja';

  const stanOdczytu = document.createElement('p');
  stanOdczytu.className = 'du-odpowiedz';
  stanOdczytu.textContent = 'Odczyt katalogu ustawień rdzenia…';

  const glowny = utworzPrzelacznik('Powiadomienia aktywne');
  glowny.kontrolka.addEventListener('change', () => {
    void zapisz(zrodlo.ustawLogiczna(KLUCZ_GLOWNY, glowny.kontrolka.checked));
  });

  const opisGlownego = document.createElement('p');
  opisGlownego.className = 'dn-pole-opis';
  opisGlownego.textContent =
    'Wyłączenie wygasza wszystkie klasy zdarzeń naraz i zachowuje ich ustawienia — ' +
    'ponowne włączenie przywraca stan poprzedni, nie domyślny. Wyłączenie nie ogranicza ' +
    'dostępu do żadnej funkcji platformy.';

  const tabela = document.createElement('table');
  tabela.className = 'dn-tabela du-powiadomienia';
  const glowa = document.createElement('thead');
  const wierszGlowy = document.createElement('tr');
  for (const naglowek of ['Klasa zdarzenia', 'Zgłaszaj', 'Kanały']) {
    const komorka = document.createElement('th');
    komorka.scope = 'col';
    komorka.textContent = naglowek;
    wierszGlowy.append(komorka);
  }
  glowa.append(wierszGlowy);
  const cialo = document.createElement('tbody');
  tabela.append(glowa, cialo);

  const oDoreczaniu = document.createElement('p');
  oDoreczaniu.className = 'dn-pole-opis du-granica';
  oDoreczaniu.textContent =
    'Wybór zapisuje się natychmiast. Samego doręczania nie ma jeszcze czym wykonać: ' +
    'silnik kolejki powiadomień jest w rdzeniu zbudowany, ale nie ma wołacza, a centrum ' +
    'powiadomień nie ma rodziny komend w kontrakcie. Zapisane tu rozstrzygnięcie czeka ' +
    'więc na wpięcie doręczeń — i będzie wtedy obowiązywać bez zmiany w tym oknie.';

  const zdanie = document.createElement('p');
  zdanie.className = 'du-odpowiedz';
  zdanie.hidden = true;

  element.append(stanOdczytu, glowny.element, opisGlownego, tabela, zdanie, oDoreczaniu);
  pokazTresc(false);

  /** Wiersze zbudowane raz; odczyt kolejny wchodzi przez `ustaw`. */
  const wiersze = new Map<KlasaZdarzenia, WierszKlasy>();

  /** Chowa treść sekcji, dopóki odczyt jej nie potwierdzi. */
  function pokazTresc(widoczna: boolean): void {
    glowny.element.hidden = !widoczna;
    opisGlownego.hidden = !widoczna;
    tabela.hidden = !widoczna;
    oDoreczaniu.hidden = !widoczna;
  }

  /** Nanosi odpowiedź czynności: milczenie przy powodzeniu, zdanie przy odmowie. */
  async function zapisz(czynnosc: Promise<string>): Promise<void> {
    const odmowa = await czynnosc;
    zdanie.hidden = odmowa === '';
    zdanie.dataset['powodzenie'] = 'false';
    zdanie.textContent = odmowa;
  }

  function nanies(stan: StanPowiadomien): void {
    if (stan.odmowa !== '') {
      stanOdczytu.hidden = false;
      stanOdczytu.dataset['powodzenie'] = 'false';
      stanOdczytu.textContent = stan.odmowa;
      pokazTresc(false);
      return;
    }

    stanOdczytu.hidden = true;
    pokazTresc(true);
    glowny.kontrolka.checked = stan.wlaczone;

    for (const klasa of stan.klasy) {
      const istniejacy = wiersze.get(klasa.klasa);
      if (istniejacy === undefined) {
        const nowy = utworzWierszKlasy(klasa, stan, zapisz, zrodlo);
        wiersze.set(klasa.klasa, nowy);
        cialo.append(nowy.element);
      } else {
        istniejacy.ustaw(klasa, stan.wlaczone);
      }
    }
    // Wygaszenie idzie po nałożeniu wierszy, bo dotyczy każdego z nich.
    for (const wiersz of wiersze.values()) wiersz.wygas(!stan.wlaczone);
  }

  function odczytaj(): void {
    stanOdczytu.hidden = false;
    stanOdczytu.dataset['powodzenie'] = 'true';
    stanOdczytu.textContent = 'Odczyt katalogu ustawień rdzenia…';
    void zrodlo.odczytaj().then(nanies);
  }

  // Zmiana z drugiego okna dolatuje zdarzeniem rdzenia; sekcja nadąża nasłuchem, nie odpytuje w tle.
  const odsubskrybuj = zrodlo.naZmiane(() => void zrodlo.odczytaj().then(nanies));
  odczytaj();

  return {
    element,
    odswiez: odczytaj,
    rozlacz: () => odsubskrybuj(),
  };
}

/** Wiersz jednej klasy zdarzenia w tabeli, niosący nazwę, ster czynności i wybór kanałów dostarczenia rdzenia. */
interface WierszKlasy {
  element: HTMLTableRowElement;
  ustaw(stan: StanKlasy, wlaczoneGlownie: boolean): void;
  /** Wygasza sterowanie wiersza, nie zmieniając ani jednej wartości. */
  wygas(wygaszony: boolean): void;
}

function utworzWierszKlasy(
  stan: StanKlasy,
  calosc: StanPowiadomien,
  zapisz: (czynnosc: Promise<string>) => Promise<void>,
  zrodlo: ZrodloPowiadomien,
): WierszKlasy {
  const element = document.createElement('tr');

  const nazwa = document.createElement('th');
  nazwa.scope = 'row';
  nazwa.className = 'du-powiadomienia__klasa';
  const etykieta = document.createElement('span');
  etykieta.className = 'du-powiadomienia__nazwa';
  etykieta.textContent = stan.definicja?.name ?? stan.klasa;
  const opis = document.createElement('span');
  opis.className = 'dn-pole-opis';
  opis.textContent = stan.definicja?.description ?? '';
  opis.hidden = opis.textContent === '';
  nazwa.append(etykieta, opis);

  const komorkaCzynnosci = document.createElement('td');
  const czynna = utworzPrzelacznik(`Zgłaszaj klasę ${etykieta.textContent}`, true);
  czynna.kontrolka.addEventListener('change', () => {
    void zapisz(zrodlo.ustawLogiczna(kluczKlasy(stan.klasa), czynna.kontrolka.checked));
  });
  komorkaCzynnosci.append(czynna.element);

  const komorkaKanalow = document.createElement('td');
  komorkaKanalow.className = 'du-powiadomienia__kanaly';
  // Centrum stoi jako stan nazwany, nie jako ster: platforma nie umie go wyłączyć z rejestru.
  const centrum = document.createElement('span');
  centrum.className = 'dn-pole-opis';
  centrum.textContent = 'centrum — zawsze';
  komorkaKanalow.append(centrum);

  const wybory = new Map<string, HTMLInputElement>();
  for (const kanal of calosc.dostepneKanaly) {
    const wybor = document.createElement('label');
    wybor.className = 'dn-wybor du-powiadomienia__kanal';
    const pole = document.createElement('input');
    pole.type = 'checkbox';
    pole.className = 'dn-check';
    pole.checked = stan.kanaly.includes(kanal.wartosc);
    pole.addEventListener('change', () => {
      const wybrane = [...wybory]
        .filter(([, kontrolka]) => kontrolka.checked)
        .map(([wartosc]) => wartosc);
      void zapisz(zrodlo.ustawKanaly(stan.klasa, wybrane));
    });
    const napis = document.createElement('span');
    napis.textContent = kanal.etykieta;
    wybor.append(pole, napis);
    wybory.set(kanal.wartosc, pole);
    komorkaKanalow.append(wybor);
  }

  element.append(nazwa, komorkaCzynnosci, komorkaKanalow);
  czynna.kontrolka.checked = stan.czynna;

  return {
    element,

    ustaw(nowy) {
      czynna.kontrolka.checked = nowy.czynna;
      for (const [wartosc, pole] of wybory) pole.checked = nowy.kanaly.includes(wartosc);
    },

    wygas(wygaszony) {
      // Wygaszenie jest widokiem, nie zapisem: wartości zostają, sterowanie przestaje przyjmować zmiany.
      element.dataset['wygaszony'] = String(wygaszony);
      czynna.kontrolka.disabled = wygaszony;
      for (const pole of wybory.values()) pole.disabled = wygaszony;
    },
  };
}

/** Przełącznik — kontrolka i element osadzany bezpośrednio w dokumencie, bo nie zawsze są tym samym elementem. */
interface Przelacznik {
  /** Element wstawiany do dokumentu: etykieta z kontrolką albo sama kontrolka. */
  element: HTMLElement;
  kontrolka: HTMLInputElement;
}

/** Przełącznik z etykietą jest jedną postacią na całą sekcję; w tabeli etykieta widoczna nie staje, bo nazwę klasy niesie nagłówek wiersza. */
function utworzPrzelacznik(etykieta: string, bezWidocznejEtykiety = false): Przelacznik {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'checkbox';
  kontrolka.className = 'dn-przelacznik';
  kontrolka.setAttribute('aria-label', etykieta);
  if (bezWidocznejEtykiety) return { element: kontrolka, kontrolka };

  const opakowanie = document.createElement('label');
  opakowanie.className = 'dn-wybor du-powiadomienia__glowny';
  const napis = document.createElement('span');
  napis.textContent = etykieta;
  opakowanie.append(kontrolka, napis);
  return { element: opakowanie, kontrolka };
}
