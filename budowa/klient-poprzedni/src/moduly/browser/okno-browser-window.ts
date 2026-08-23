import { utworzPanelRodzin } from './panel-rodzin';
import { sekcjePrzegladania } from './sekcje-rodzin';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import { opisOdmowy } from '../../komponenty/odmowa';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { utworzCzynnosciPaskaDolnego } from './czynnosci-paska-dolnego';
import {
  KLASY_DYMKA,
  KODY_OKIEN,
  KODY_PANELI,
  OBJASNIENIA,
  STANY_PUSTE,
  WYZWALACZE,
} from './etykiety-browser';
import { utworzFormularzNawigacji } from './formularz-nawigacji';
import { utworzNarzedziaInspekcyjne } from './narzedzia-inspekcyjne';
import { utworzPanelWyodrebnien } from './panel-wyodrebnien';
import { utworzPasekDolny } from './pasek-dolny';
import { utworzPasekZaznaczenia } from './pasek-zaznaczenia';
import { utworzPodgladStrony } from './podglad-strony';
import { utworzStanOkna } from './stan-okna';
import type { StanPrzegladania } from './stan-przegladania';
import { utworzWarstweAdnotacji } from './warstwa-adnotacji';
import type { RozszerzenieModulu } from './warstwy-widocznosci';

/**
 * Browser Window — okno wiodące modułu przeglądarki.
 *
 * Jedna odpowiedzialność: rama okna. Nie buduje ani jednego elementu treści —
 * składa gotowe części: formularz nawigacji, podgląd strony, pasek zaznaczenia,
 * pasek dolny i panel wyodrębnień.
 *
 * Trzy czynności Operatora: nawigacja do strony (formularz), przewijanie
 * i zaznaczenie fragmentu (podgląd). Pierwsza idzie do rdzenia; przewijanie ma
 * komendę, której okno jeszcze nie wywołuje, a zaznaczenie dzieje się wyłącznie
 * w kliencie i komendy nie potrzebuje.
 *
 * Odmowa rdzenia jest treścią okna, nie jego awarią. Każda komenda obszaru
 * `browser.*` może odmówić — okno pokazuje wtedy powód odmowy i mówi wprost,
 * co odczytało naprawdę, zamiast pustego prostokąta.
 */
export interface OknoBrowserWindow {
  element: HTMLElement;
  /** Nanosi stan modułu na wszystkie części okna. */
  odswiez(): void;
  /** Odczytuje stan okna komunikacji (`window.state.get`). */
  wczytajStan(): Promise<void>;
  /** Panele warstwy czwartej okna — moduł rejestruje je w warstwach widoczności. */
  panele: readonly RozszerzenieModulu[];
  /** Odpina warstwę adnotacji — moduł kończy pracę. */
  rozlacz(): void;
}

export function utworzOknoBrowserWindow(
  stan: StanPrzegladania,
  naNotatke: (fragment: string) => void,
): OknoBrowserWindow {
  const okno = utworzStanOkna(STANY_PUSTE.przegladarka);
  const odpowiedz = utworzWierszOdpowiedzi();
  const powiedz = (tresc: string, powodzenie: boolean): void =>
    odpowiedz.pokaz(tresc, powodzenie);

  // Czynności strony powstają raz i obsługują oba paski: dolny i pływający.
  const czynnosciStrony = utworzCzynnosciPaskaDolnego(stan, powiedz);

  const podglad = utworzPodgladStrony((fragment) => stan.ustawZaznaczenie(fragment));
  const zaznaczenie = utworzPasekZaznaczenia(stan, {
    naNotatke,
    naTlumaczenie: () => czynnosciStrony.doTlumaczenia(),
    powiedz,
  });
  const wyodrebnienia = utworzPanelWyodrebnien(stan);
  const nawigacja = utworzFormularzNawigacji(stan, {
    przewin: (kierunek) => podglad.przewin(kierunek),
    pokazWiersz: (numer) => podglad.pokazWiersz(numer),
    powiedz,
  });

  // Warstwa adnotacji powstaje przed paskiem dolnym, bo pasek jest jej
  // przełącznikiem — a przełącznik nie może wskazywać czegoś, czego jeszcze nie
  // ma. Zwrotne `naZmianeTrybu` domyka pętlę: tryb zamknięty pływającym paskiem
  // gasi przycisk paska dolnego.
  const adnotacja = utworzWarstweAdnotacji(stan, {
    powiedz,
    naZmianeTrybu: (wlaczona) => dolny.ustawTrybAdnotacji(wlaczona),
  });

  const dolny = utworzPasekDolny(czynnosciStrony, {
    naPodzial: (wlaczony) => {
      cialo.dataset['podzial'] = wlaczony ? 'tak' : 'nie';
    },
    naCzytnik: (wlaczony) => podglad.ustawTrybCzytnika(wlaczony),
    naAdnotacje: (wlaczony) => adnotacja.ustaw(wlaczony),
    naWyodrebnienie: (rodzaj) => {
      const tresc = wyodrebnienia.wyodrebnij(rodzaj);
      powiedz(
        tresc === ''
          ? 'Wyodrębnienie nic nie dało — migawka nie niesie tej treści.'
          : `Wyodrębniono ${tresc.length} znaków treści strony.`,
        tresc !== '',
      );
    },
    powiedz,
  });

  const stanOkna = przycisk('Odczytaj stan okna', 'dn-btn dn-btn--sm dn-btn--zarys');
  const opisStanu = document.createElement('p');
  opisStanu.className = 'mb-stan-okna';

  // Wskaźnik obecności Wykonawcy. Mówi o migawce, bo to ona jest wspólnym
  // widokiem: model widzi dokładnie tę treść i dokładnie z tej chwili, a nie
  // stronę taką, jaka jest teraz w sieci. Stan niesie napis, nie sama barwa.
  const obecnosc = document.createElement('span');
  obecnosc.className = 'dn-plakietka dn-plakietka--informacja mb-obecnosc';

  // Scena jest jedynym kontenerem pozycjonującym warstwę. Płótno adnotacji
  // jest rodzeństwem podglądu strony, nie jego dzieckiem, więc oznaczanie
  // strony nie wstrzykuje w jej DOM ani jednego węzła.
  const scena = document.createElement('div');
  scena.className = 'mb-scena';
  scena.append(podglad.element, adnotacja.element);

  const cialo = document.createElement('div');
  cialo.className = 'mb-przegladarka__cialo';
  cialo.dataset['podzial'] = 'nie';
  cialo.append(scena, wyodrebnienia.element);

  // Panel inspekcji staje nad obszarem renderowania, a nie obok niego: to, co
  // narzędzia pokazują, dotyczy strony leżącej pod nimi.
  const inspekcja = utworzNarzedziaInspekcyjne(stan, powiedz);

  // Rodziny prowadzone przez rdzeń: karty, grupy kart, przestrzenie robocze,
  // zakładki, przewinięcie, zrzut i narzędzia inspekcyjne CDP. Panel stoi pod
  // obszarem renderowania, bo dotyczy strony leżącej wyżej — i nie zasłania jej.
  const rodziny = utworzPanelRodzin(sekcjePrzegladania(stan));

  okno.tresc.append(
    nawigacja.element,
    inspekcja.element,
    cialo,
    zaznaczenie.element,
    rodziny.element,
    odpowiedz.element,
  );

  const naglowek = utworzNaglowekOkna({
    tytul: 'Browser Window',
    kontrolki: [obecnosc, utworzDymekObjasnienia(OBJASNIENIA.obecnoscWykonawcy, KLASY_DYMKA), stanOkna, opisStanu],
    klasa: 'mb-okno__naglowek',
  });

  const element = document.createElement('section');
  element.className = 'mb-okno mb-okno--wiodace';
  element.dataset['okno'] = KODY_OKIEN.przegladarka;
  element.setAttribute('aria-label', 'Browser Window — podgląd współdzielony strony');
  element.append(naglowek, okno.element, dolny.element);

  async function wczytajStan(): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') {
      opisStanu.textContent = stan.powod();
      return;
    }
    opisStanu.textContent = 'Odczyt stanu okna w toku…';
    const wynik = await stan.okna.stan(idOkna);
    if (!wynik.udany || wynik.wynik === undefined) {
      opisStanu.textContent = opisOdmowy('Odczyt stanu okna', wynik.blad?.code, wynik.blad?.message);
      return;
    }
    const tresc = wynik.wynik;
    opisStanu.textContent =
      `Stan procesu: ${tresc.processStatus} · wiadomości: ${tresc.messageCount} · ` +
      `strumień: ${tresc.streaming ? 'trwa' : 'nie trwa'}`;
  }

  stanOkna.addEventListener('click', () => void wczytajStan());
  // Ponowienie sięga po jedno i drugie: stan okna czyta `window.state.get`,
  // a treść strony `browser.snapshot.get`, i to ta druga odmowa stawia okno
  // w stanie błędu. Ponowienie odczytujące tylko stan zostawiałoby Operatora
  // z komunikatem, którego nie da się zdjąć czynnością, którą mu podano.
  okno.ustawPonowienie(() => {
    void wczytajStan();
    void stan.zaciagnijMigawke();
  });

  return {
    element,

    odswiez() {
      podglad.odswiez(stan.migawka());
      zaznaczenie.odswiez();
      wyodrebnienia.odswiez();
      obecnosc.textContent = zdanieObecnosci(stan);
      // Trzy stany w jednej kolejności: najpierw czekanie (obie drogi odczytu),
      // potem odmowa, na końcu pustka. Odwrotna kolejność kazałaby Operatorowi
      // czytać „jest pusto" w chwili, w której odczyt jeszcze trwa.
      if (stan.faza() === 'odczyt') {
        okno.ladowanie('Rdzeń ustala okno przeglądarki tej sesji…');
        return;
      }
      if (stan.fazaMigawki() === 'odczyt') {
        okno.ladowanie(stan.powodMigawki());
        return;
      }
      if (stan.faza() === 'blad') {
        okno.blad(stan.powod());
        return;
      }
      if (stan.migawka() === null) {
        // Pustka i odmowa to dwie różne rzeczy. „Rdzeń nie ma jeszcze migawki
        // tego okna" jest stanem pustym; nieudany odczyt treści strony jest
        // odmową i ma się nią przedstawić, bo Operator inaczej zareaguje na
        // jedno, a inaczej na drugie.
        if (stan.fazaMigawki() === 'blad') okno.blad(stan.powodMigawki());
        else okno.puste(STANY_PUSTE.przegladarka.tytul, zdaniePustki(stan));
        return;
      }
      okno.gotowe();
    },

    wczytajStan,

    panele: [
      {
        kod: KODY_PANELI.narzedziaInspekcyjne,
        nazwa: WYZWALACZE.inspekcja,
        warstwa: 4,
        element: inspekcja.element,
        // Skrót z załącznika opracowania modułu. Powłoka gospodarza potrafi
        // przejąć tę kombinację przed stroną — wtedy zostaje wyzwalacz paska
        // kontekstu w trybie administracyjnym.
        skrot: { klawisz: 'i', zShift: true },
      },
    ],

    rozlacz: adnotacja.rozlacz,
  };
}

/**
 * Zdanie stanu pustego: czym okno jest, a zaraz po tym — co dokładnie mówi
 * o nim rdzeń w tej chwili.
 *
 * Obie części są konieczne. Sam opis okna nie mówi, czy przeszkodą jest brak
 * sesji, brak okna przeglądarki czy tylko brak przejścia; sam powód z rdzenia
 * nie mówi, czym okno jest ani jak je zapełnić.
 *
 * Bez okna przeglądarki nie ma o co pytać o migawkę, więc zdanie o migawce
 * opisywałoby wtedy skutek zamiast przyczyny — Operator ma dostać powód braku
 * okna, bo to on rozstrzyga, co da się zrobić dalej.
 */
/**
 * Zdanie wskaźnika obecności Wykonawcy.
 *
 * Wspólny podgląd jest wspólny przez migawkę: model dostaje jej treść, więc
 * wskaźnik mówi o migawce i o chwili jej pobrania, a nie o „obecności" bez
 * pokrycia. Bez migawki mówi wprost, że model nie ma na czym pracować.
 */
function zdanieObecnosci(stan: StanPrzegladania): string {
  const migawka = stan.migawka();
  if (migawka === null) return 'Wykonawca bez widoku — okno nie ma jeszcze migawki';
  return `Wykonawca widzi migawkę z ${new Date(migawka.capturedAt).toLocaleTimeString('pl-PL')}`;
}

function zdaniePustki(stan: StanPrzegladania): string {
  const zRdzenia = stan.idOkna() === '' ? stan.powod() : stan.powodMigawki();
  return zRdzenia === ''
    ? STANY_PUSTE.przegladarka.opis
    : `${STANY_PUSTE.przegladarka.opis} ${zRdzenia}`;
}
