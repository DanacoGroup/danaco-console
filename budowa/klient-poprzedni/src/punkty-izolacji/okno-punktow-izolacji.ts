import './punkty-izolacji.css';

import { ChangeKind, EventType } from '../../../shared/contract';
import { utworzRameOkna } from '../komponenty/rama-okna';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { PANELE_IZOLACJI, REJESTR_OBSZAROW, type ObszarIzolacji } from './obszary';
import { utworzStanWarstwy } from './stan-warstwy';
import { utworzStanZasiegu } from './stan-zasiegu';
import { utworzSterowanieWarstwa } from './sterowanie-warstwa';

/**
 * Okno Punktów Izolacji — okno zarządca katalogu rdzenia (`punkty-izolacji`,
 * kategoria `konfiguracja`), otwierane z listwy Ustawienia obok Okna
 * Konfiguracji, Dostępów i Modeli.
 *
 * Dwie warstwy, żadnej drugiej ramy. Obudowa modalna (`<dialog>`, tło, Esc,
 * `dn-modal`) jest tym samym idiomem co w oknie konfiguracji, dostępów
 * i modeli. Wewnątrz ciała modalu stoi panel zbudowany przez
 * `komponenty/rama-okna.ts` (tytuł, plakietka roli, przeznaczenie, pas
 * narzędzi, ciało) — biblioteczny komponent współdzielony z oknami
 * operacyjnymi modułów, a nie drugi, przepisany od nowa.
 *
 * Plik zna cztery rzeczy: obudowę modalną, wybór warstwy izolacji, podział na
 * trzy panele i miejsce montażu obszarów w każdym z nich. Nie zna ani jednego
 * klucza izolacji — te dostarczają pliki `obszar-*.ts` przez kontrakt
 * `ObszarIzolacji` z `obszary.ts`.
 *
 * Trzy panele, nie zakładki (rozdz. 6.1 Modelu konfiguracji): selektor zasięgu
 * po lewej, macierz izolacji pośrodku, profil i podgląd polityki efektywnej po
 * prawej. Panele stoją obok siebie i są widoczne naraz, bo mówią o jednej
 * rzeczy w trzech ujęciach: gdzie reguła obowiązuje, co ustawia i co z tego
 * wynika po dziedziczeniu. Rozdzielone na zakładki kazałyby Operatorowi
 * pamiętać wybór zasięgu z jednej zakładki, przestawiając przełącznik
 * w drugiej.
 *
 * Warstwa izolacji (`default` — bazowa platformy, `session` — nakładka karty
 * sesji wygasająca z jej zamknięciem) oraz zasięg (poziom i byt, `stan-zasiegu.ts`)
 * rozstrzygają, który zapis czyta i pisze każdy obszar. Oba są stanem wspólnym
 * okna: warstwa stoi w pasie narzędzi ramy, zasięg w panelu lewym, a ich zmiana
 * odświeża wszystkie obszary naraz — inaczej jeden panel pokazywałby wartości
 * zasięgu, którego w selektorze już nie ma.
 *
 * Okno otwiera się natychmiast, przed jakąkolwiek odpowiedzią rdzenia; komendy
 * `isolation.*` wołają dopiero obszary, każdy własnym odczytem.
 *
 * Rdzeń ogłasza dwa zdarzenia — `isolation.profile.changed`
 * i `isolation.policy.changed` — i tutaj są one słuchane raz, dla całego okna,
 * a nie osobno w każdym obszarze. Jedna subskrypcja odświeża wszystkie sześć
 * obszarów; osobne byłyby rozjeżdżającymi się odczytami tej samej zmiany. Bez
 * nich punkt przestawiony z drugiego okna albo ręką asystenta zostawiałby
 * w oknie wartość nieświeżą, bez śladu i bez odświeżenia.
 *
 * Odświeżenie nie jest jedyną odpowiedzią: sama zmiana wartości pod palcami
 * Operatora byłaby posunięciem niewidzialnym, więc nad obszarem staje ślad —
 * co się zmieniło i czego to dotyczyło.
 *
 * Ślad nie powie, czyja ręka. Oba zdarzenia idą bez pól
 * `actor`/`actorClientId`, choć rdzeń wypełnia je w innych kopertach
 * (`core/sprawca.go`), więc ślad mówi „nie wiadomo, czyja ręka" — napis
 * pewniejszy niż dowód byłby gorszy od jego braku.
 */
export interface OknoPunktowIzolacji {
  /** Element `<dialog>` osadzany w dokumencie. */
  element: HTMLDialogElement;
  /** Otwiera okno i montuje obszar czynny. */
  otworz(): void;
  /** Zamyka okno; stan i subskrypcje zostają. */
  zamknij(): void;
  /** Odłącza subskrypcje obszaru czynnego i usuwa okno z dokumentu. */
  rozlacz(): void;
}

export function utworzOknoPunktowIzolacji(kanal: Kanal): OknoPunktowIzolacji {
  const rama = utworzRameOkna({
    kod: 'punkty-izolacji',
    tytul: 'Punkty izolacji',
    rola: 'zarządca',
    modul: 'Punkty izolacji',
    przeznaczenie:
      'Konfiguracja jedenastu punktów izolacji rdzenia — kontekst, zakres techniczny, profile, ' +
      'poziom zapisu, oś rozstrzygania i podgląd polityki efektywnej, na wybranej warstwie.',
    przedrostek: 'pi',
  });

  const warstwa = utworzStanWarstwy();
  const zasieg = utworzStanZasiegu();
  const sterowanieWarstwa = utworzSterowanieWarstwa(kanal, warstwa);

  // Ślad ostatniej zmiany przyszłej z rdzenia. Stoi nad obszarem, bo dotyczy
  // każdego obszaru naraz; ukryty, dopóki nic się nie zmieniło — pusty wiersz
  // „brak zmian" byłby szumem, a nie odpowiedzią.
  const slad = document.createElement('p');
  slad.className = 'pi-slad';
  slad.setAttribute('role', 'status');
  slad.hidden = true;

  rama.narzedzia.append(sterowanieWarstwa.element);

  // Sześć obszarów montowanych naraz, po dwa na panel. Egzemplarz każdego
  // powstaje raz na życie okna: przemontowywanie ich przy każdym odświeżeniu
  // gubiłoby treść wpisaną w formularze profilu i pola punktu widzenia.
  const obszary: ObszarIzolacji[] = [];
  const kolumny = PANELE_IZOLACJI.map((panel) => {
    const kolumna = zbudujKolumne(panel.tytul, panel.opis, panel.kod);
    for (const wpis of REJESTR_OBSZAROW.filter((w) => w.panel === panel.kod)) {
      const obszar = wpis.utworz({ kanal, warstwa, zasieg });
      obszary.push(obszar);
      kolumna.append(zbudujMiejsceObszaru(wpis.nazwa, wpis.opis, obszar.element));
    }
    return kolumna;
  });

  const panele = document.createElement('div');
  panele.className = 'pi-panele';
  panele.append(...kolumny);

  rama.cialo.append(slad, panele);

  /** Odczyt wszystkich obszarów od nowa — po otwarciu okna i po każdej zmianie wspólnego stanu. */
  function odswiezWszystkie(): void {
    for (const obszar of obszary) obszar.odswiez();
  }

  // Zmiana warstwy i zmiana zasięgu dotyczą każdego obszaru naraz, więc czytają
  // od nowa wszystkie — inaczej macierz pokazywałaby wartości zasięgu albo
  // warstwy, której w selektorze już nie ma.
  warstwa.naZmiane(() => odswiezWszystkie());
  zasieg.naZmiane(() => odswiezWszystkie());

  /**
   * Nanosi zmianę zgłoszoną przez rdzeń: ślad na wierzchu i ponowny odczyt
   * obszaru czynnego.
   *
   * Odczyt idzie tylko wtedy, gdy okno jest otwarte. Zamknięte okno nie ma co
   * odświeżać — `otworz()` montuje obszar od nowa i sam woła `odswiez()`.
   * Ślad zostaje mimo to zapisany, żeby Operator zobaczył go przy wejściu.
   */
  function nanieszZmiane(zdanie: string): void {
    slad.hidden = false;
    slad.textContent =
      `${zdanie} Odczytano od nowa. Sprawcy nie znamy: to zdarzenie nie niesie pola ` +
      '„actor" — nie wiadomo, czy punkt przestawił Operator z innego urządzenia, ' +
      'czy asystent.';
    if (!element.open) return;
    odswiezWszystkie();
  }

  const odsubskrybowania: Odsubskrybuj[] = [
    kanal.naZdarzenie(EventType.IsolationProfileChanged, (tresc) => {
      const nazwa = tresc.profileId === '' ? 'bez nazwy' : tresc.profileId;
      nanieszZmiane(`Profil izolacji „${nazwa}" — ${slowoZmiany(tresc.change)}.`);
    }),
    kanal.naZdarzenie(EventType.IsolationPolicyChanged, (tresc) => {
      const okno = tresc.windowId === '' ? 'nie wskazano okna' : `okno ${tresc.windowId}`;
      nanieszZmiane(`Polityka izolacji (${okno}) — ${slowoZmiany(tresc.change)}.`);
    }),
  ];

  const element = document.createElement('dialog');
  element.className = 'dn-modal pi-okno';
  element.setAttribute('aria-label', 'Okno punktów izolacji');

  const cialo = document.createElement('div');
  cialo.className = 'dn-modal-cialo pi-okno__cialo';
  cialo.append(rama.element);

  const zamknij = document.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn dn-btn--atrament';
  zamknij.textContent = 'Zamknij';
  zamknij.addEventListener('click', () => element.close());

  const stopka = document.createElement('footer');
  stopka.className = 'dn-modal-stopka pi-okno__stopka';
  stopka.append(zamknij);

  element.append(cialo, stopka);

  return {
    element,

    otworz() {
      if (!element.isConnected) document.body.append(element);
      if (!element.open) element.showModal();
      sterowanieWarstwa.odswiez();
      odswiezWszystkie();
    },

    zamknij: () => element.close(),

    rozlacz() {
      for (const odsubskrybuj of odsubskrybowania.splice(0)) odsubskrybuj();
      for (const obszar of obszary.splice(0)) obszar.zamknij?.();
      element.remove();
    },
  };
}

/** Kolumna panelu: tytuł, zdanie o przeznaczeniu, miejsce na obszary. */
function zbudujKolumne(tytul: string, opis: string, kod: string): HTMLElement {
  const naglowek = document.createElement('h3');
  naglowek.className = 'pi-panel__tytul';
  naglowek.textContent = tytul;

  const zdanie = document.createElement('p');
  zdanie.className = 'pi-panel__opis';
  zdanie.textContent = opis;

  const element = document.createElement('section');
  element.className = 'pi-panel';
  element.dataset['panel'] = kod;
  element.setAttribute('aria-label', tytul);
  element.append(naglowek, zdanie);
  return element;
}

/** Obszar w panelu: podpis obszaru nad jego treścią — panel mieści po dwa. */
function zbudujMiejsceObszaru(nazwa: string, opis: string, tresc: HTMLElement): HTMLElement {
  const podpis = document.createElement('h4');
  podpis.className = 'pi-obszar__tytul';
  podpis.textContent = nazwa;

  const zdanie = document.createElement('p');
  zdanie.className = 'pi-obszar__opis';
  zdanie.textContent = opis;

  const element = document.createElement('div');
  element.className = 'pi-obszar';
  element.append(podpis, zdanie, tresc);
  return element;
}

/** Rodzaj zmiany po polsku — `ChangeKind` jest kontraktem, nie napisem dla Operatora. */
function slowoZmiany(zmiana: ChangeKind): string {
  switch (zmiana) {
    case ChangeKind.Created:
      return 'utworzony';
    case ChangeKind.Updated:
      return 'zmieniony';
    case ChangeKind.Deleted:
      return 'usunięty';
    default:
      // Rodzaj spoza wykazu nie jest powodem do milczenia: pokazujemy go
      // dosłownie, tak jak przyszedł.
      return `zmiana rodzaju „${String(zmiana)}"`;
  }
}

