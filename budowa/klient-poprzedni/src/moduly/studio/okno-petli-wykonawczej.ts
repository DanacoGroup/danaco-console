import './petla-okno.css';

import {
  StudioOperationScope,
  StudioTaskState,
  type StudioAgentConflict,
  type StudioChain,
  type StudioDocumentTask,
  type StudioOperation,
} from '../../../../shared/contract';
import { utworzWidokObsady } from './agenci-obsada';
import { utworzOknoStudio } from './okno-studio';
import { utworzStanPetli } from './petla-stan';
import { utworzWarsztatLancucha } from './petla-warsztat-lancucha';
import { utworzTrybWsadowy } from './petla-wsad';
import type { ZrodloPetli } from './petla-zrodlo';
import type { StanStudio } from './stan-studio';
import { utworzWykazZadan, przycisk } from './zadania-wykaz';

/**
 * Execution Loop Window — okno pętli wykonawczej modułu Studio.
 *
 * ── POWIERZCHNIA NALEŻY DO DOKUMENTU ────────────────────────────────────────
 * Opracowanie modułu stawia to okno jako kolumnę sąsiadującą z oknem rozmowy.
 * Rozstrzygnięcie Właściciela, podjęte dwukrotnie i obowiązujące w całym module,
 * jest mocniejsze: panele nie zajmują stałych kolumn — wchodzą na żądanie
 * i schodzą, gdy nie są używane. Okno pętli wchodzi więc znacznikiem przebiegu
 * i schodzi po zamknięciu, a stała kolumna jest TRYBEM DO WYBORU Operatora,
 * jawnym, odwracalnym i pamiętanym — nigdy postacią domyślną.
 *
 * Samoczynne otwarcie zostaje, bo jest wprost w opracowaniu: zlecenie
 * wielozadaniowe otwiera okno w chwili uruchomienia. Otwarcie samoczynne nie
 * jest tu sprzeczne z zasadą powierzchni — okno wchodzi wtedy, gdy naprawdę jest
 * używane, i to jest dokładnie ta reguła.
 *
 * ── Cztery widoki, jedna kolejka ────────────────────────────────────────────
 * Kolejka zadań, obsada wykonawców, warsztat łańcucha i tryb wsadowy patrzą na
 * ten sam stan. Krok łańcucha i dokument wsadu są zadaniami tej samej kolejki, bo
 * kontrakt mówi wprost: przebieg łańcucha i wsadu prowadzi pętla wykonawcza okna.
 * Druga maszyneria obok byłaby drugą prawdą o tym samym przebiegu.
 *
 * ── Czego to okno nie robi ──────────────────────────────────────────────────
 * Nie przestawia nastaw pętli. Nastawy — czy pętla jest czynna, ilu wykonawców
 * naraz, co przy spięciu — idą zasięgami rodziny `config.*` i komendą
 * `studio.agents.settings.set`. Okno je CZYTA i pokazuje wprost, a gdy pętla jest
 * wyłączona, mówi to zdaniem, zamiast pozwalać naciskać przycisk, który i tak
 * odmówi.
 */

/** Wywołania okna sięgające do sąsiednich okien modułu. */
export interface CzynnosciSasiedzkie {
  /** Otwiera wynik zadania w oknie pracy z dokumentem. */
  pokazWynikZadania(zadanie: StudioDocumentTask): void;
  /** Przenosi wskazanie na fragment dokumentu. */
  pokazFragment(od: number, do_: number): void;
}

export interface OknoPetliWykonawczej {
  /** Element okna osadzany w układzie modułu. */
  element: HTMLElement;
  /** Znacznik przebiegu — wchodzi do paska kontekstu i otwiera okno. */
  znacznik: HTMLElement;
  /** Odczytuje rozkład, obsadę, spięcia, łańcuchy i nastawy. */
  wczytaj(): Promise<void>;
  /** Przerysowuje okno ze stanu. */
  odswiez(): void;
  /** Rozkłada zlecenie dokumentowe na zadania i otwiera okno. */
  rozlozZlecenie(zlecenie: string): Promise<void>;
  /** Odłącza subskrypcję postępu przebiegu. */
  zamknij(): void;
}

export function utworzOknoPetliWykonawczej(
  stanStudio: StanStudio,
  zrodlo: ZrodloPetli,
  sasiedzkie: CzynnosciSasiedzkie,
): OknoPetliWykonawczej {
  const stan = utworzStanPetli();

  const okno = utworzOknoStudio({
    kod: 'studio.execution-loop',
    tytul: 'Pętla wykonawcza',
    rola: 'pomocnicze',
    objasnienie:
      'Rozkłada zlecenie dokumentowe na zadania i prowadzi ich kolejkę — komendami ' +
      'studio.plan.create, .get, .task.update, .run i .stop. Pętla wykonawcza i praca ' +
      'kilku wykonawców naraz są narzędziami NIEURUCHAMIANYMI NA STARCIE: włącza je ' +
      'nastawa Operatora (studio.agents.settings.set), a dopóki stoi na „nie", ' +
      'uruchomienie wraca odmową nazywającą brak nastawy.',
  });

  // ── Znacznik przebiegu ────────────────────────────────────────────────────
  const znacznik = document.createElement('button');
  znacznik.type = 'button';
  znacznik.className = 'dn-btn dn-btn--zarys dn-btn--sm petla-znacznik';
  znacznik.addEventListener('click', () => {
    stan.ustawOtwarte(!stan.otwarte());
    if (stan.otwarte()) void wczytajWszystko();
  });

  // ── Pasek sterowania przebiegiem ──────────────────────────────────────────
  const nastawaZdanie = document.createElement('p');
  nastawaZdanie.className = 'dn-tekst-3 petla-nastawa';

  const wskazniki = document.createElement('p');
  wskazniki.className = 'dn-tekst-3 petla-wskazniki';

  const puscGuzik = przycisk('Puść pętlę', () => void puscPetle(false));
  const probaGuzik = przycisk('Sprawdź bez ruszania dokumentu', () => void puscPetle(true));
  const stopGuzik = przycisk('Przerwij', () => void zatrzymaj());
  const kolumnaGuzik = przycisk('Stała kolumna', () => {
    stan.ustawStalaKolumne(!stan.stalaKolumna());
  });
  const zamknijGuzik = przycisk('Zwiń okno', () => stan.ustawOtwarte(false));

  const sterowanie = document.createElement('div');
  sterowanie.className = 'petla-sterowanie';
  sterowanie.append(puscGuzik, probaGuzik, stopGuzik, kolumnaGuzik, zamknijGuzik);

  // ── Zawężenie kolejki wedle stanu zadania ─────────────────────────────────
  const filtrStanu = document.createElement('select');
  filtrStanu.className = 'dn-wybor petla-filtr';
  filtrStanu.setAttribute('aria-label', 'Zawężenie kolejki do stanu zadania');
  for (const [wartosc, napis] of [
    ['', 'wszystkie stany'],
    [StudioTaskState.Pending, 'czekające'],
    [StudioTaskState.Running, 'w realizacji'],
    [StudioTaskState.Done, 'zakończone'],
    [StudioTaskState.Failed, 'nieudane'],
    [StudioTaskState.Blocked, 'wstrzymane'],
    [StudioTaskState.Skipped, 'pominięte'],
  ] as const) {
    const pozycja = document.createElement('option');
    pozycja.value = wartosc;
    pozycja.textContent = napis;
    filtrStanu.append(pozycja);
  }
  filtrStanu.addEventListener('change', () => void wczytajRozklad());

  const wykazZadan = utworzWykazZadan({
    pomin: (idZadania) => void przestaw(idZadania, StudioTaskState.Skipped),
    ponow: (idZadania) => void przestaw(idZadania, StudioTaskState.Pending),
    pokazWynik: (zadanie) => sasiedzkie.pokazWynikZadania(zadanie),
  });

  const widokObsady = utworzWidokObsady({
    przyjmijBrzmienie: (spiecie) => void przyjmijBrzmienie(spiecie),
    pokazFragment: (od, do_) => sasiedzkie.pokazFragment(od, do_),
  });

  const warsztat = utworzWarsztatLancucha({
    zapisz: () => void zapiszLancuch(),
    uruchom: (idLancucha) => void uruchomLancuch(idLancucha),
    wczytajDoZmiany: (lancuch: StudioChain) => {
      stan.ustawZmieniany(lancuch.id);
      stan.ustawNazweSkladanego(lancuch.name);
      stan.ustawSkladane(lancuch.steps);
    },
    ustawKroki: (kroki) => stan.ustawSkladane(kroki),
    ustawNazwe: (nazwa) => stan.ustawNazweSkladanego(nazwa),
  });

  const wsad = utworzTrybWsadowy({
    puscWsad: (dokumenty, idAkcji) => void puscWsad(dokumenty, idAkcji),
    przerwij: () => stan.ustawPrzerwanieWsadu(true),
  });

  let operacjeWlasne: readonly StudioOperation[] = [];

  const zakladki = utworzZakladki([
    ['kolejka', 'Kolejka zadań', wykazZadan.element],
    ['agenci', 'Kto nad czym pracuje', widokObsady.element],
    ['lancuch', 'Warsztat łańcucha', warsztat.element],
    ['wsad', 'Tryb wsadowy', wsad.element],
  ]);

  okno.stan.tresc.append(nastawaZdanie, wskazniki, sterowanie, filtrStanu, zakladki.element);

  // Znacznik przebiegu stoi POZA oknem, bo okno bywa zwinięte, a znacznik jest
  // jedyną drogą do jego rozwinięcia. Pas obejmuje oba: jest kotwicą układu dla
  // nakładki (`position: absolute` potrzebuje przodka ustawionego) i miejscem,
  // w którym znacznik zostaje widoczny przy zwiniętym oknie.
  const pas = document.createElement('div');
  pas.className = 'petla-pas';
  pas.append(znacznik, okno.element);

  // ── Postęp przebiegu z rdzenia ────────────────────────────────────────────
  // Okno odświeża się ZDARZENIEM, nie odpytywaniem w pętli. Rozkład o dziesięciu
  // zadaniach odpytywany co sekundę znaczyłby dziesięć wywołań na nic w każdej
  // sekundzie, w której nic się nie stało.
  const odsubskrybujPostep = zrodlo.naPostep((tresc) => {
    const rozklad = stan.rozklad();
    if (rozklad === null || tresc.runId !== rozklad.id) return;
    void wczytajRozklad();
  });

  async function wczytajNastawy(): Promise<void> {
    const dokument = stanStudio.dokument();
    const wynik = await zrodlo.nastawy(
      dokument === null ? {} : { documentId: dokument.id, windowId: stanStudio.idOkna() },
    );
    if (wynik.udany && wynik.wynik !== undefined) {
      stan.ustawNastawy(wynik.wynik.settings);
      return;
    }
    // Odczyt nieudany zostawia nastawy jako „nie wiem" — a nie jako
    // „wyłączone". Zgadnięcie „wyłączone" byłoby twierdzeniem o stanie rdzenia,
    // którego okno nie sprawdziło.
    stan.ustawNastawy(null);
    okno.stan.blad(
      `Nastaw pętli nie udało się odczytać: ${wynik.blad?.message ?? 'bez powodu podanego przez rdzeń'}`,
    );
  }

  async function wczytajRozklad(): Promise<void> {
    const dokument = stanStudio.dokument();
    if (dokument === null) {
      stan.ustawRozklad(null);
      return;
    }
    const zawezenie = filtrStanu.value === '' ? undefined : (filtrStanu.value as StudioTaskState);
    const wynik = await zrodlo.rozklad({
      documentId: dokument.id,
      ...(zawezenie === undefined ? {} : { state: zawezenie }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      // Dokument bez ani jednego rozkładu wraca odmową `not_found` — i to jest
      // odpowiedź prawdziwa, nie awaria. Okno pokazuje wtedy stan pusty, a nie
      // błąd: „rozkładu nie ma" i „odczyt się nie udał" to dwie różne rzeczy.
      stan.ustawRozklad(null);
      stan.ustawRozklady([]);
      okno.stan.puste(
        'Rozkładu nie ma',
        'Ten dokument nie ma jeszcze rozłożonego zlecenia. Wydaj zlecenie w oknie rozmowy ' +
          'albo rozłóż je wprost — pętla ułoży zadania i pokaże je tutaj.',
      );
      return;
    }
    stan.ustawRozklad(wynik.wynik.plan);
    stan.ustawRozklady(wynik.wynik.plans ?? [wynik.wynik.plan]);
    okno.stan.gotowe();
  }

  async function wczytajObsade(): Promise<void> {
    const dokument = stanStudio.dokument();
    if (dokument === null) return;
    const [obsada, spiecia] = await Promise.all([
      zrodlo.obsada({ documentId: dokument.id }),
      zrodlo.spiecia({ documentId: dokument.id }),
    ]);
    if (obsada.udany && obsada.wynik !== undefined) stan.ustawObsade(obsada.wynik.slots);
    if (spiecia.udany && spiecia.wynik !== undefined) stan.ustawSpiecia(spiecia.wynik.conflicts);
  }

  async function wczytajLancuchy(): Promise<void> {
    const [lancuchy, wlasne] = await Promise.all([
      zrodlo.lancuchy({}),
      zrodlo.operacjeWlasne({}),
    ]);
    if (lancuchy.udany && lancuchy.wynik !== undefined) stan.ustawLancuchy(lancuchy.wynik.chains);
    if (wlasne.udany && wlasne.wynik !== undefined) {
      operacjeWlasne = wlasne.wynik.operations.filter((operacja) => !operacja.builtin);
    }
  }

  async function wczytajWszystko(): Promise<void> {
    okno.stan.ladowanie('Odczyt rozkładu, obsady i łańcuchów z rdzenia');
    await Promise.all([wczytajNastawy(), wczytajRozklad(), wczytajObsade(), wczytajLancuchy()]);
  }

  async function przestaw(idZadania: string, docelowy: StudioTaskState): Promise<void> {
    const wynik = await zrodlo.przestawZadanie({ taskId: idZadania, state: docelowy });
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.stan.blad(
        `Zadania nie przestawiono: ${wynik.blad?.message ?? 'rdzeń nie podał powodu'}`,
      );
      return;
    }
    stan.ustawRozklad(wynik.wynik.plan);
  }

  async function puscPetle(probny: boolean): Promise<void> {
    const rozklad = stan.rozklad();
    if (rozklad === null) return;
    stan.ustawOdmowePrzebiegu(null);
    const wynik = await zrodlo.puscPetle({
      planId: rozklad.id,
      ...(probny ? { dryRun: true } : {}),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.stan.blad(`Pętli nie puszczono: ${wynik.blad?.message ?? 'rdzeń nie podał powodu'}`);
      return;
    }
    stan.ustawRozklad(wynik.wynik.plan);
    // Pętla wyłączona nastawą wraca `started: false` wraz z powodem. Powód
    // WCHODZI DO OKNA wprost — po to rdzeń go nazywa. Milczące nic nie zmieniło
    // by się od udanego przebiegu bez pracy, a to jest wzorzec szkody.
    if (!wynik.wynik.started) {
      stan.ustawOdmowePrzebiegu(
        wynik.wynik.refusalReason ?? 'rdzeń odmówił uruchomienia bez podania powodu',
      );
      return;
    }
    stan.ustawOdmowePrzebiegu(null);
    const bilans = wynik.wynik.balance;
    if (bilans !== undefined && bilans.skippedCount > 0) {
      // Pominięcia przebiegu — także te z blokad fragmentu — pokazujemy nad
      // kolejką, bo dotyczą przebiegu w całości, a nie jednego zadania.
      okno.stan.blad(zdanieOPominieciach(bilans.note, bilans.skippedCount));
    } else {
      okno.stan.gotowe();
    }
  }

  async function zatrzymaj(): Promise<void> {
    const rozklad = stan.rozklad();
    if (rozklad === null) return;
    stan.ustawPrzerwanieWsadu(true);
    const wynik = await zrodlo.zatrzymaj({
      planId: rozklad.id,
      reason: 'przerwanie z okna pętli wykonawczej',
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.stan.blad(`Pętli nie zatrzymano: ${wynik.blad?.message ?? 'rdzeń nie podał powodu'}`);
      return;
    }
    stan.ustawRozklad(wynik.wynik.plan);
  }

  async function zapiszLancuch(): Promise<void> {
    const kroki = stan.skladane();
    const nazwa = stan.nazwaSkladanego().trim();
    if (kroki.length === 0 || nazwa === '') return;
    const zmieniany = stan.zmieniany();
    const wynik = await zrodlo.zapiszLancuch({
      name: nazwa,
      steps: [...kroki],
      ...(zmieniany === '' ? {} : { chainId: zmieniany }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.stan.blad(
        `Łańcucha nie zapisano: ${wynik.blad?.message ?? 'rdzeń nie podał powodu'}`,
      );
      return;
    }
    stan.ustawZmieniany('');
    stan.ustawSkladane([]);
    stan.ustawNazweSkladanego('');
    await wczytajLancuchy();
  }

  /**
   * Uruchomienie łańcucha: rdzeń rejestruje przebieg, a kroki wchodzą do
   * kolejki pętli jako zadania.
   *
   * Kontrakt mówi wprost, że przebieg łańcucha prowadzi pętla wykonawcza okna —
   * więc okno robi dwie rzeczy w tej kolejności: woła `studio.chain.run`, żeby
   * przebieg powstał i był policzony po stronie rdzenia, a potem rozkłada kroki
   * na zadania (`studio.plan.create`), żeby każdy krok miał swój stan, swojego
   * wykonawcę i swój powód, gdy się nie uda.
   */
  async function uruchomLancuch(idLancucha: string): Promise<void> {
    const dokument = stanStudio.dokument();
    if (dokument === null) {
      okno.stan.blad('Łańcucha nie ma na czym uruchomić — żaden dokument nie jest czynny.');
      return;
    }
    const lancuch = stan.lancuchy().find((wpis) => wpis.id === idLancucha);
    if (lancuch === undefined) return;

    const zaznaczenie = stanStudio.zaznaczenie();
    const przebieg = await zrodlo.uruchomLancuch({
      windowId: stanStudio.idOkna(),
      documentId: dokument.id,
      chainId: idLancucha,
      scope:
        zaznaczenie === null ? StudioOperationScope.Document : StudioOperationScope.Selection,
      ...(zaznaczenie === null
        ? {}
        : { selectionStart: zaznaczenie.poczatek, selectionEnd: zaznaczenie.koniec }),
    });
    if (!przebieg.udany) {
      okno.stan.blad(
        `Łańcucha nie uruchomiono: ${przebieg.blad?.message ?? 'rdzeń nie podał powodu'}`,
      );
      return;
    }

    // Każdy krok łańcucha wchodzi jako zadanie: rodzaj `custom`, bo krok jest
    // czynnością Operatora z katalogu, a nie jedną z siedmiu faz zlecenia.
    const zadania = lancuch.steps.map((krok, numer) => ({
      id: '',
      planId: '',
      documentId: dokument.id,
      order: numer + 1,
      kind: 'custom' as const,
      title: `Krok ${numer + 1}: ${krok.actionId}`,
      instruction: krok.actionId,
      state: StudioTaskState.Pending,
    }));
    const rozklad = await zrodlo.rozlozZlecenie({
      documentId: dokument.id,
      order: `łańcuch „${lancuch.name}" — ${lancuch.steps.length} kroków`,
      tasks: zadania,
      ...(zaznaczenie === null
        ? {}
        : { rangeStart: zaznaczenie.poczatek, rangeEnd: zaznaczenie.koniec }),
    });
    if (!rozklad.udany || rozklad.wynik === undefined) {
      okno.stan.blad(
        `Kroków łańcucha nie wniesiono do kolejki: ${rozklad.blad?.message ?? 'bez powodu'}`,
      );
      return;
    }
    stan.ustawRozklad(rozklad.wynik.plan);
    stan.ustawOtwarte(true);
    zakladki.wybierz('kolejka');
    await puscPetle(false);
  }

  /**
   * Wsad: dokument po dokumencie, z chwilą na „przerwij" między nimi.
   *
   * Zatrzymanie NIE wycofuje niczego. Dokumenty już przetworzone zostają
   * przetworzone; pozostałe dostają stan „przerwana" wraz z powodem, żeby
   * Operator wiedział, gdzie wsad się skończył, a nie zgadywał.
   */
  async function puscWsad(dokumenty: readonly string[], idAkcji: string): Promise<void> {
    stan.ustawPrzerwanieWsadu(false);
    stan.ustawWsad(
      dokumenty.map((idDokumentu) => ({ idDokumentu, stan: 'oczekuje', powod: '' })),
    );
    stan.ustawOtwarte(true);
    zakladki.wybierz('wsad');

    for (const idDokumentu of dokumenty) {
      if (stan.wsadPrzerwany()) {
        stan.przestawPozycjeWsadu(
          idDokumentu,
          'przerwana',
          'wsad przerwany decyzją Operatora przed tym dokumentem; dokumenty przetworzone ' +
            'wcześniej zostały przetworzone i nic się nie wycofuje',
        );
        continue;
      }
      stan.przestawPozycjeWsadu(idDokumentu, 'w-biegu', '');
      const wynik = await zrodlo.uruchomWsad({
        windowId: stanStudio.idOkna(),
        documentIds: [idDokumentu],
        actionId: idAkcji,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        stan.przestawPozycjeWsadu(
          idDokumentu,
          'odrzucona',
          wynik.blad?.message ?? 'rdzeń odmówił bez podania powodu',
        );
        continue;
      }
      const odrzucenie = wynik.wynik.rejected?.[0];
      if (odrzucenie !== undefined) {
        stan.przestawPozycjeWsadu(idDokumentu, 'odrzucona', odrzucenie.reason);
        continue;
      }
      if (wynik.wynik.accepted === 0) {
        // Rdzeń nie przyjął dokumentu i nie podał odrzucenia. Nie udajemy
        // powodzenia: „przyjęto 0" bez wykazu jest bilansem niepełnym i tak to
        // nazywamy.
        stan.przestawPozycjeWsadu(
          idDokumentu,
          'odrzucona',
          'rdzeń nie przyjął dokumentu i nie podał odrzucenia — bilans wsadu jest niepełny',
        );
        continue;
      }
      stan.przestawPozycjeWsadu(idDokumentu, 'udana', '');
    }
  }

  /**
   * Przyjęcie odłożonego brzmienia po spięciu.
   *
   * Brzmienie odłożone wchodzi do dokumentu decyzją Operatora, a nie samo.
   * Droga jest ta sama, którą wchodzi każda zmiana: przez rozkład i pętlę, więc
   * przyjęcie ma swój wiersz w kolejce i swój ślad. Rdzeń nie ma dziś komendy
   * „przyjmij odłożone brzmienie" — to jest brak wypisany w sprawozdaniu, a nie
   * dziura zaklejona wywołaniem, którego nie ma.
   */
  async function przyjmijBrzmienie(spiecie: StudioAgentConflict): Promise<void> {
    if (spiecie.deferredText === undefined || spiecie.deferredText === '') return;
    const rozklad = await zrodlo.rozlozZlecenie({
      documentId: spiecie.documentId,
      order:
        `przyjmij brzmienie odłożone po spięciu na znakach ` +
        `${spiecie.rangeStart}–${spiecie.rangeEnd}`,
      tasks: [
        {
          id: '',
          planId: '',
          documentId: spiecie.documentId,
          order: 1,
          kind: 'draft' as const,
          title: 'Wniesienie odłożonego brzmienia',
          instruction: spiecie.deferredText,
          state: StudioTaskState.Pending,
          rangeStart: spiecie.rangeStart,
          rangeEnd: spiecie.rangeEnd,
        },
      ],
      rangeStart: spiecie.rangeStart,
      rangeEnd: spiecie.rangeEnd,
    });
    if (!rozklad.udany || rozklad.wynik === undefined) {
      okno.stan.blad(
        `Odłożonego brzmienia nie wniesiono: ${rozklad.blad?.message ?? 'bez powodu'}`,
      );
      return;
    }
    stan.ustawRozklad(rozklad.wynik.plan);
    zakladki.wybierz('kolejka');
    await puscPetle(false);
  }

  function odswiez(): void {
    const rozklad = stan.rozklad();
    const nastawy = stan.nastawy();
    const zadania = rozklad?.tasks ?? [];
    const wBiegu = zadania.filter((zadanie) => zadanie.state === StudioTaskState.Running).length;
    const nieudane = zadania.filter((zadanie) => zadanie.state === StudioTaskState.Failed).length;

    znacznik.textContent =
      rozklad === null
        ? '⇄ pętla'
        : `⇄ pętla · ${zadania.length} zadań${wBiegu > 0 ? ` · ${wBiegu} w biegu` : ''}`;
    znacznik.setAttribute('aria-expanded', stan.otwarte() ? 'true' : 'false');

    okno.element.hidden = !stan.otwarte();
    okno.element.dataset['tryb'] = stan.stalaKolumna() ? 'kolumna' : 'nakladka';
    kolumnaGuzik.textContent = stan.stalaKolumna()
      ? 'Wróć do nakładki na żądanie'
      : 'Stała kolumna (tryb do wyboru)';

    if (nastawy === null) {
      nastawaZdanie.textContent =
        'Nastaw pętli jeszcze nie odczytano. Okno nie twierdzi, że pętla jest wyłączona — ' +
        'twierdzi, że nie pytało.';
    } else {
      nastawaZdanie.textContent =
        `Pętla wykonawcza: ${nastawy.executionLoopEnabled ? 'CZYNNA' : 'wyłączona'} · ` +
        `praca wielu wykonawców: ${nastawy.multiAgentEnabled ? 'CZYNNA' : 'wyłączona'}` +
        (nastawy.maxConcurrentAgents !== undefined
          ? ` · najwyżej ${nastawy.maxConcurrentAgents} naraz`
          : '') +
        (nastawy.conflictPolicy !== undefined
          ? ` · przy spięciu: ${nastawy.conflictPolicy}`
          : '') +
        '. Oba narzędzia są domyślnie wyłączone i włącza je nastawa Operatora ' +
        '(studio.agents.settings.set), nie to okno.';
    }

    const odmowa = stan.odmowaPrzebiegu();
    wskazniki.textContent =
      rozklad === null
        ? ''
        : `Przebieg: ${zadania.length} zadań · ${rozklad.iterations ?? 0} obiegów · ` +
          `${rozklad.iterationsWithoutProgress ?? 0} bez postępu · ${nieudane} nieudanych · ` +
          `stan rozkładu: ${rozklad.state}` +
          (rozklad.stopReason !== undefined ? ` · ${rozklad.stopReason}` : '');
    if (odmowa !== null) {
      wskazniki.textContent = `Rdzeń ODMÓWIŁ uruchomienia: ${odmowa}`;
      wskazniki.dataset['odmowa'] = 'tak';
    } else {
      delete wskazniki.dataset['odmowa'];
    }

    // Przycisk, który i tak odmówi, jest nieosiągalny — a powód stoi zdaniem
    // wyżej. Ukrycie powodu i zostawienie czynnego przycisku byłoby zaproszeniem
    // do odmowy bez wyjaśnienia.
    const mozeRuszyc = rozklad !== null && nastawy !== null && nastawy.executionLoopEnabled;
    puscGuzik.disabled = !mozeRuszyc;
    probaGuzik.disabled = rozklad === null;
    stopGuzik.disabled = rozklad === null;

    wykazZadan.odswiez(rozklad);
    widokObsady.odswiez(stan.obsada(), stan.spiecia());
    warsztat.odswiez(
      stan.skladane(),
      stan.nazwaSkladanego(),
      stan.lancuchy(),
      operacjeWlasne,
    );
    wsad.odswiez(stan.wsad(), stan.wsadPrzerwany(), operacjeWlasne);
  }

  stan.obserwuj(odswiez);
  odswiez();

  return {
    element: pas,
    znacznik,

    async wczytaj() {
      // Odczyt przy wczytaniu modułu ogranicza się do nastaw: okno jest zwinięte,
      // a odczytywanie rozkładu, obsady i łańcuchów dla okna, którego Operator
      // nie otworzył, byłoby czterema wywołaniami na nic. Nastawy są wyjątkiem,
      // bo od nich zależy treść znacznika przebiegu, który JEST widoczny.
      await wczytajNastawy();
    },

    odswiez,

    async rozlozZlecenie(zlecenie) {
      const dokument = stanStudio.dokument();
      if (dokument === null) {
        okno.stan.blad('Zlecenia nie ma na czym rozłożyć — żaden dokument nie jest czynny.');
        return;
      }
      const wynik = await zrodlo.rozlozZlecenie({ documentId: dokument.id, order: zlecenie });
      if (!wynik.udany || wynik.wynik === undefined) {
        okno.stan.blad(
          `Zlecenia nie rozłożono: ${wynik.blad?.message ?? 'rdzeń nie podał powodu'}`,
        );
        return;
      }
      stan.ustawRozklad(wynik.wynik.plan);
      // Zlecenie WIELOZADANIOWE otwiera okno samoczynnie — tak każe opracowanie
      // modułu. Zlecenie jednozadaniowe okna nie otwiera: jedno zadanie widać
      // w pasku kontekstu i nie warto mu oddawać powierzchni dokumentu.
      if ((wynik.wynik.plan.tasks?.length ?? 0) > 1) stan.ustawOtwarte(true);
      await wczytajNastawy();
    },

    zamknij() {
      odsubskrybujPostep();
    },
  };
}

/** Zdanie o pominięciach przebiegu — liczba plus nota rdzenia. */
function zdanieOPominieciach(nota: string | undefined, ile: number): string {
  const wstep = `Przebieg zostawił ${ile} pominięć — pętla ich NIE przemilcza`;
  return nota === undefined || nota === '' ? `${wstep}.` : `${wstep}: ${nota}`;
}

/** Zakładki widoków okna — cztery widoki, jedna powierzchnia. */
interface ZakladkiOkna {
  element: HTMLElement;
  wybierz(kod: string): void;
}

function utworzZakladki(
  widoki: readonly (readonly [string, string, HTMLElement])[],
): ZakladkiOkna {
  const pasek = document.createElement('div');
  pasek.className = 'dn-zakladki petla-zakladki';
  pasek.setAttribute('role', 'tablist');

  const ciala = document.createElement('div');
  ciala.className = 'petla-zakladki__ciala';

  const guziki = new Map<string, HTMLButtonElement>();

  function wybierz(kod: string): void {
    for (const [wpisKod, , cialo] of widoki) {
      cialo.hidden = wpisKod !== kod;
      const guzik = guziki.get(wpisKod);
      if (guzik !== undefined) guzik.setAttribute('aria-selected', String(wpisKod === kod));
    }
  }

  for (const [kod, napis, cialo] of widoki) {
    const guzik = document.createElement('button');
    guzik.type = 'button';
    guzik.className = 'dn-zakladka';
    guzik.setAttribute('role', 'tab');
    guzik.textContent = napis;
    guzik.addEventListener('click', () => wybierz(kod));
    guziki.set(kod, guzik);
    pasek.append(guzik);
    ciala.append(cialo);
  }

  const element = document.createElement('div');
  element.className = 'petla-widoki';
  element.append(pasek, ciala);

  const pierwszy = widoki[0];
  if (pierwszy !== undefined) wybierz(pierwszy[0]);

  return { element, wybierz };
}
