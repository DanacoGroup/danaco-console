import './dziennik-kontrola.css';

import {
  ConfigScope,
  StudioAgentConflictPolicy,
  StudioLockScope,
  type StudioAgentSlot,
  type StudioFragmentLock,
} from '../../../../shared/contract';
import { poleLogiczne, poleTekstowe, poleWyboru } from '../../modele/kontrolki-formularza';
import {
  BLOKADA_STOI_W_RDZENIU,
  BLOKADA_ZDEJMUJE_OPERATOR,
  blokadaOpisz,
  blokadaOpiszOdmoweZajecia,
  blokadaOpiszWykaz,
  blokadaOpiszZajecie,
  blokadaZakresPoprawny,
} from './blokada-fragmentow';
import type { ZrodloKontroliStudio } from './zrodlo-kontroli-studio';

/**
 * Panel blokad fragmentów, zajęć wykonawców i nastaw pracy kilku agentów naraz.
 *
 * ── Trzy rzeczy w jednym miejscu, bo dotyczą jednego pytania ────────────────
 * „Czego modelowi nie wolno tknąć i kto teraz pisze po którym akapicie."
 *
 *   — **blokada** jest trwała i skierowana przeciw modelowi; zdejmuje ją
 *     wyłącznie Operator, a zasięg obejmujący także Operatora jest osobnym,
 *     jawnym ustawieniem, nie zachowaniem domyślnym;
 *   — **zajęcie fragmentu** jest czasowe i skierowane przeciw drugiemu
 *     wykonawcy; odmowa nazywa wykonawcę i czas, bo cicha odmowa kazałaby
 *     Operatorowi zgadywać, dlaczego fragment nie drgnął;
 *   — **nastawy** rozstrzygają, czy pętla wykonawcza i praca wielu agentów
 *     w ogóle stoją, ilu wykonawców pracuje naraz i co się dzieje przy spięciu.
 *     Oba narzędzia są domyślnie wyłączone i włącza je Operator, nie okno.
 *
 * Panel woła rdzeń sam, bo skutkiem każdej z tych czynności jest wykaz albo
 * odmowa nazwana — jedno i drugie jest treścią dla Operatora, nie wartością
 * pośrednią przekazywaną wyżej.
 */

/** Czym panel pyta okno o dokument i zaznaczenie. */
export interface KontekstBlokad {
  idDokumentu(): string;
  zaznaczenie(): { poczatek: number; koniec: number } | null;
}

export interface BlokadaPanel {
  element: HTMLElement;
  /** Odczytuje blokady dokumentu. */
  odswiez(): Promise<void>;
  /** Blokady ostatnio odczytane — oknu do oznaczenia fragmentów w treści. */
  blokady(): readonly StudioFragmentLock[];
  przestawWidocznosc(): void;
  widoczny(): boolean;
}

export function utworzBlokadaPanel(
  zrodlo: ZrodloKontroliStudio,
  kontekst: KontekstBlokad,
): BlokadaPanel {
  let otwarty = false;
  let wykazBlokad: readonly StudioFragmentLock[] = [];

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'dn-pole-opis ms-kontrola__odpowiedz';
  odpowiedz.dataset['czynnosc'] = 'odpowiedz';

  function powiedz(tresc: string, udana: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.dataset['udana'] = udana ? 'tak' : 'nie';
  }

  function dokument(): string {
    const kod = kontekst.idDokumentu();
    if (kod === '') {
      powiedz(
        'Nie ma dokumentu czynnego. Blokady i zajęcia fragmentów dotyczą jednego dokumentu — ' +
          'otwórz go albo załóż nowy.',
        false,
      );
    }
    return kod;
  }

  /** Zaznaczenie nadające się na fragment albo odpowiedź nazywająca brak. */
  function fragment(): { poczatek: number; koniec: number } | null {
    const zakres = kontekst.zaznaczenie();
    if (!blokadaZakresPoprawny(zakres)) {
      powiedz(
        'Nic nie jest zaznaczone. Blokada i zajęcie dotyczą FRAGMENTU, nie całego dokumentu — ' +
          'zaznacz w treści to, co ma zostać dosłownie.',
        false,
      );
      return null;
    }
    return zakres;
  }

  /* ── Blokady ────────────────────────────────────────────────────────────── */

  const zasadaRdzenia = document.createElement('p');
  zasadaRdzenia.className = 'dn-pole-opis ms-kontrola__zasada';
  zasadaRdzenia.textContent = BLOKADA_STOI_W_RDZENIU;

  const zasadaZdejmowania = document.createElement('p');
  zasadaZdejmowania.className = 'dn-pole-opis ms-kontrola__zasada';
  zasadaZdejmowania.textContent = BLOKADA_ZDEJMUJE_OPERATOR;

  const nazwaBlokady = poleTekstowe({
    etykieta: 'Nazwa blokady',
    podpowiedz: 'na przykład „podstawa prawna"',
  });
  const powodBlokady = poleTekstowe({
    etykieta: 'Powód blokady',
    podpowiedz: 'na przykład „podstawa prawna — nie zmieniać"',
    opis: 'Powód wraca w odmowie, którą model dostaje przy próbie zmiany — warto, by mówił za siebie.',
  });
  const zasiegBlokady = poleWyboru(
    {
      etykieta: 'Kogo blokada dotyczy',
      opis:
        'Blokada jest skierowana przeciw MODELOWI, nie przeciw właścicielowi dokumentu. Objęcie ' +
        'nią Operatora jest osobnym, jawnym ustawieniem, a nie zachowaniem domyślnym.',
    },
    [
      { wartosc: StudioLockScope.Model, etykieta: 'Wyłącznie model' },
      { wartosc: StudioLockScope.Everyone, etykieta: 'Model i Operator' },
    ],
  );

  const zalozBlokade = przyciskPanelu('Załóż blokadę na zaznaczeniu', 'blokada-zaloz');
  zalozBlokade.addEventListener('click', () => {
    void zaloz();
  });

  const podsumowanieBlokad = document.createElement('p');
  podsumowanieBlokad.className = 'dn-pole-opis';

  const wykazBlokadElement = document.createElement('ul');
  wykazBlokadElement.className = 'ms-kontrola__wykaz';

  /* ── Zajęcia fragmentów ─────────────────────────────────────────────────── */

  const idWykonawcy = poleTekstowe({
    etykieta: 'Wykonawca zajmujący fragment',
    podpowiedz: 'identyfikator eksperta',
    opis:
      'Puste znaczy wykonawcę nienazwanego. Zajęcie bez nazwy wykonawcy da się założyć, ale odmowa ' +
      'drugiemu nie nazwie wtedy nikogo — dlatego warto je podać.',
  });
  const nazwaWykonawcy = poleTekstowe({
    etykieta: 'Nazwa wykonawcy widoczna dla Operatora',
    podpowiedz: 'nazwa, którą Operator rozpozna',
  });
  const ttl = poleTekstowe({
    etykieta: 'Po ilu sekundach zajęcie samo wygasa',
    podpowiedz: 'na przykład 300',
    opis:
      'Wygasanie jest tu zabezpieczeniem, nie szczegółem: wykonawca ubity w pół pracy nie ma ' +
      'trzymać fragmentu na zawsze.',
  });

  const zajmij = przyciskPanelu('Zajmij zaznaczony fragment', 'zajecie-zajmij');
  zajmij.addEventListener('click', () => {
    void zajmijFragment();
  });
  const zwolnij = przyciskPanelu('Zwolnij wszystkie zajęcia wykonawcy', 'zajecie-zwolnij');
  zwolnij.addEventListener('click', () => {
    void zwolnijZajecia();
  });

  const wykazZajec = document.createElement('ul');
  wykazZajec.className = 'ms-kontrola__wykaz';

  /* ── Nastawy pracy wielu wykonawców ─────────────────────────────────────── */

  const petla = poleLogiczne({
    etykieta: 'Pętla wykonawcza czynna',
    opis: 'Domyślnie wyłączona. Włącza ją Operator jawnym, odwracalnym ustawieniem.',
  });
  const wielu = poleLogiczne({
    etykieta: 'Praca wielu wykonawców naraz czynna',
    opis: 'Domyślnie wyłączona — dwóch agentów nad jednym pismem to decyzja Operatora.',
  });
  const ilu = poleTekstowe({
    etykieta: 'Ilu wykonawców może pracować naraz',
    podpowiedz: 'na przykład 3',
  });
  const przySpieciu = poleWyboru(
    {
      etykieta: 'Co robić przy spięciu o ten sam fragment',
      opis:
        'Odmowa nazywa fragment i wykonawcę; odłożenie zachowuje brzmienie drugiego zamiast ' +
        'nadpisać pierwsze; blokada wykonawcy trzyma fragment na czas pracy.',
    },
    [
      { wartosc: StudioAgentConflictPolicy.Refuse, etykieta: 'Odmowa drugiemu' },
      { wartosc: StudioAgentConflictPolicy.Queue, etykieta: 'Odłożenie zmiany drugiego' },
      { wartosc: StudioAgentConflictPolicy.FragmentLock, etykieta: 'Blokada fragmentu na czas pracy' },
    ],
  );
  const wymagajZajecia = poleLogiczne({
    etykieta: 'Wykonawca musi zająć fragment, zanim go zmieni',
  });
  const zasiegNastaw = poleWyboru(
    {
      etykieta: 'Zasięg nastaw',
      opis: 'Wartości idą zasięgami rodziny config, nie osobnym magazynem Studia.',
    },
    [
      { wartosc: ConfigScope.Window, etykieta: 'Okno' },
      { wartosc: ConfigScope.Session, etykieta: 'Karta sesji' },
      { wartosc: ConfigScope.Module, etykieta: 'Moduł Studio' },
      { wartosc: ConfigScope.Global, etykieta: 'Cała maszyna' },
    ],
  );

  const zapiszNastawy = przyciskPanelu('Zapisz nastawy wykonawców', 'nastawy-wykonawcow');
  zapiszNastawy.addEventListener('click', () => {
    void ustawNastawy();
  });

  const zdanieNastaw = document.createElement('p');
  zdanieNastaw.className = 'dn-pole-opis';

  /* ── Złożenie panelu ────────────────────────────────────────────────────── */

  const element = document.createElement('section');
  element.className = 'ms-kontrola';
  element.dataset['panel'] = 'blokady-i-wykonawcy';
  element.hidden = true;
  element.setAttribute(
    'aria-label',
    'Blokady fragmentów, zajęcia wykonawców i nastawy pracy wielu agentów',
  );
  element.append(
    czescPanelu('Blokada fragmentu — czego model nie tknie', [
      zasadaRdzenia,
      zasadaZdejmowania,
      nazwaBlokady.element,
      powodBlokady.element,
      zasiegBlokady.element,
      pasPrzyciskow([zalozBlokade]),
      podsumowanieBlokad,
      wykazBlokadElement,
    ]),
    czescPanelu('Zajęcie fragmentu przez wykonawcę', [
      idWykonawcy.element,
      nazwaWykonawcy.element,
      ttl.element,
      pasPrzyciskow([zajmij, zwolnij]),
      wykazZajec,
    ]),
    czescPanelu('Nastawy pracy wielu wykonawców', [
      petla.element,
      wielu.element,
      ilu.element,
      przySpieciu.element,
      wymagajZajecia.element,
      zasiegNastaw.element,
      pasPrzyciskow([zapiszNastawy]),
      zdanieNastaw,
    ]),
    odpowiedz,
  );

  /* ── Czynności ──────────────────────────────────────────────────────────── */

  function liczba(pole: HTMLInputElement): number | undefined {
    const wpisane = pole.value.trim();
    if (wpisane === '') return undefined;
    const wartosc = Number(wpisane);
    return Number.isFinite(wartosc) && wartosc >= 0 ? wartosc : undefined;
  }

  async function zaloz(): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const zakres = fragment();
    if (zakres === null) return;
    const nazwa = nazwaBlokady.kontrolka.value.trim();
    if (nazwa === '') {
      powiedz(
        'Blokada musi mieć nazwę: to ona wraca w odmowie, którą model dostaje przy próbie zmiany. ' +
          'Blokada bez nazwy nie powiedziałaby, co zatrzymało czynność.',
        false,
      );
      return;
    }
    const wynik = await zrodlo.kontrolaBlokadaZaloz(
      kod,
      zakres,
      nazwa,
      powodBlokady.kontrolka.value.trim(),
      zasiegBlokady.kontrolka.value as StudioLockScope,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Blokady nie założono: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    wykazBlokad = wynik.wynik.locks;
    przerysujBlokady();
    powiedz(
      `Blokada założona: ${blokadaOpisz(wynik.wynik.lock)}. ${BLOKADA_STOI_W_RDZENIU}`,
      true,
    );
  }

  async function zdejmij(blokada: StudioFragmentLock): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const wynik = await zrodlo.kontrolaBlokadaZdejmij(kod, blokada.id);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Blokady nie zdjęto: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    wykazBlokad = wynik.wynik.locks;
    przerysujBlokady();
    powiedz(
      wynik.wynik.removed
        ? `Blokada „${blokada.name}" zdjęta przez Operatora.`
        : `Blokady „${blokada.name}" rdzeń NIE zdjął. ${BLOKADA_ZDEJMUJE_OPERATOR}`,
      wynik.wynik.removed,
    );
  }

  async function odswiez(): Promise<void> {
    const kod = kontekst.idDokumentu();
    if (kod === '') return;
    const wynik = await zrodlo.kontrolaBlokadyWykaz(kod);
    if (!wynik.udany || wynik.wynik === undefined) {
      podsumowanieBlokad.textContent = `Blokad nie udało się odczytać: ${wynik.blad?.message ?? ''}`;
      wykazBlokadElement.replaceChildren();
      return;
    }
    wykazBlokad = wynik.wynik.locks;
    przerysujBlokady();
  }

  function przerysujBlokady(): void {
    podsumowanieBlokad.textContent = blokadaOpiszWykaz(wykazBlokad);
    wykazBlokadElement.replaceChildren(...wykazBlokad.map(wierszBlokady));
  }

  async function zajmijFragment(): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const zakres = fragment();
    if (zakres === null) return;
    const wynik = await zrodlo.kontrolaZajmijFragment(
      kod,
      zakres,
      {
        idWykonawcy: idWykonawcy.kontrolka.value.trim(),
        nazwaWykonawcy: nazwaWykonawcy.kontrolka.value.trim(),
      },
      liczba(ttl.kontrolka) ?? -1,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Zajęcia fragmentu nie założono: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    const tresc = wynik.wynik;
    if (!tresc.claimed) {
      // Odmowa zajęcia jest odpowiedzią UDANĄ. Bez tego rozróżnienia okno
      // pokazałoby powodzenie tam, gdzie fragmentu nie zajęto.
      powiedz(blokadaOpiszOdmoweZajecia(tresc.heldBy, tresc.refusalReason), false);
      pokazZajecia(tresc.heldBy === undefined ? [] : [tresc.heldBy]);
      return;
    }
    powiedz(
      tresc.slot === undefined
        ? 'Fragment zajęty, ale rdzeń nie oddał zajęcia — nie ma czego pokazać w wykazie.'
        : `Fragment zajęty: ${blokadaOpiszZajecie(tresc.slot)}`,
      true,
    );
    pokazZajecia(tresc.slot === undefined ? [] : [tresc.slot]);
  }

  async function zwolnijZajecia(): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const wynik = await zrodlo.kontrolaZwolnijFragment(
      kod,
      '',
      { idWykonawcy: idWykonawcy.kontrolka.value.trim() },
      true,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Zajęć nie zwolniono: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    powiedz(`Zwolnionych zajęć: ${wynik.wynik.released}.`, true);
    pokazZajecia(wynik.wynik.slots);
  }

  function pokazZajecia(zajecia: readonly StudioAgentSlot[]): void {
    wykazZajec.replaceChildren(
      ...zajecia.map((zajecie) => {
        const pozycja = document.createElement('li');
        pozycja.dataset['zajecie'] = 'czynne';
        pozycja.dataset['wykonawca'] = zajecie.actor.agentId ?? '';
        pozycja.textContent = blokadaOpiszZajecie(zajecie);
        return pozycja;
      }),
    );
  }

  async function ustawNastawy(): Promise<void> {
    const wynik = await zrodlo.kontrolaNastawyWykonawcow({
      zasieg: zasiegNastaw.kontrolka.value as ConfigScope,
      petlaCzynna: petla.kontrolka.checked,
      wieluCzynnych: wielu.kontrolka.checked,
      ilu: liczba(ilu.kontrolka),
      przySpieciu: przySpieciu.kontrolka.value as StudioAgentConflictPolicy,
      wymagajZajecia: wymagajZajecia.kontrolka.checked,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Nastaw wykonawców nie zapisano: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    const nastawy = wynik.wynik.settings;
    zdanieNastaw.textContent =
      `Pętla wykonawcza: ${nastawy.executionLoopEnabled ? 'czynna' : 'wyłączona'} · ` +
      `praca wielu naraz: ${nastawy.multiAgentEnabled ? 'czynna' : 'wyłączona'} · ` +
      `wykonawców naraz: ${nastawy.maxConcurrentAgents ?? 'liczby rdzeń nie podał'} · ` +
      `przy spięciu: ${nastawy.conflictPolicy ?? 'nastawy rdzeń nie podał'} · ` +
      `zajęcie fragmentu wymagane: ${nastawy.requireFragmentClaim === true ? 'tak' : 'nie'} · ` +
      `zasięg wartości: ${nastawy.scope ?? 'nie podany'}`;
    powiedz('Nastawy pracy wykonawców zapisane.', true);
  }

  function wierszBlokady(blokada: StudioFragmentLock): HTMLElement {
    const opis = document.createElement('p');
    opis.className = 'dn-pole-opis';
    opis.textContent = blokadaOpisz(blokada);

    const zdjecie = przyciskPanelu('Zdejmij blokadę', 'blokada-zdejmij');
    zdjecie.title = BLOKADA_ZDEJMUJE_OPERATOR;
    zdjecie.addEventListener('click', () => {
      void zdejmij(blokada);
    });

    const pozycja = document.createElement('li');
    pozycja.dataset['blokada'] = blokada.id;
    pozycja.dataset['zasieg'] = blokada.scope;
    pozycja.append(opis, zdjecie);
    return pozycja;
  }

  return {
    element,
    odswiez,
    blokady: () => wykazBlokad,

    przestawWidocznosc() {
      otwarty = !otwarty;
      element.hidden = !otwarty;
      if (otwarty) void odswiez();
    },

    widoczny: () => otwarty,
  };
}

/** Przycisk panelu wraz z kodem czynności do sprawdzianu. */
function przyciskPanelu(nazwa: string, kod: string): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  przycisk.textContent = nazwa;
  przycisk.dataset['czynnosc'] = kod;
  return przycisk;
}

/** Pas przycisków — jeden rząd czynności. */
function pasPrzyciskow(przyciski: readonly HTMLElement[]): HTMLElement {
  const pas = document.createElement('div');
  pas.className = 'ms-kontrola__pas';
  pas.append(...przyciski);
  return pas;
}

/** Część panelu wraz z jej tytułem. */
function czescPanelu(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-kontrola__tytul';
  naglowek.textContent = tytul;

  const sekcja = document.createElement('section');
  sekcja.className = 'ms-kontrola__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}
