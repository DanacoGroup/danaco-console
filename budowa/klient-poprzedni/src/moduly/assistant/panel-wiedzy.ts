import { KnowledgeScope, type KnowledgeHit } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  pole,
  pozycjaWykazu,
  przelacznik,
  przyciskAkcji,
  utworzWierszOdpowiedzi,
  wiersz,
  wybor,
  wykaz,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { BRAKI, zglosBrak } from './braki-kontraktu';
import { ODCZYTY, POZYCJE_ZAKRESU_WIEDZY, PUSTE } from './etykiety-assistant';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAssistant } from './stan-assistant';
import type { ZrodloPamieci } from './zrodlo-pamieci';

/**
 * Pamięć semantyczna i baza wiedzy profilu — dwie zakładki nad jedną parą
 * komend, więc jeden plik.
 *
 * `knowledge.search` odnajduje fragmenty PO ZNACZENIU i oddaje je wraz ze
 * źródłem, żeby dało się je zacytować zamiast streszczać. `knowledge.index`
 * buduje wskaźnik, z którego to wyszukiwanie korzysta. Rozdzielenie ich na dwa
 * pliki dałoby dwa miejsca mówiące o jednym wskaźniku.
 *
 * Zakładka rozstrzyga wyłącznie o tym, co jest w panelu widoczne: pamięć
 * semantyczna pokazuje szukanie, baza wiedzy — budowę wskaźnika. Wspólny jest
 * zakres (`KnowledgeScope`), bo wskaźnik jest jeden i szuka się w tym, co się
 * zindeksowało.
 *
 * Wgrywania dokumentów tu nie ma i nie powinno być: pliki wchodzą do platformy
 * przez moduł Library, a `knowledge.index` obejmuje wskaźnikiem to, co
 * w bibliotece już leży. Drugie wejście dla plików znaczyłoby dwa repozytoria.
 */
export interface PanelWiedzy {
  /** Obszar zakładki „Pamięć semantyczna" — szukanie po znaczeniu. */
  szukanie: HTMLElement;
  /** Obszar zakładki „Baza wiedzy" — budowa wskaźnika znaczenia. */
  wskaznik: HTMLElement;
}

export function utworzPanelWiedzy(stan: StanAssistant, zrodlo: ZrodloPamieci): PanelWiedzy {
  const okno: StanOkna = utworzStanOkna();
  const odpowiedz: WierszOdpowiedzi = utworzWierszOdpowiedzi();

  const zapytanie = pole('Pytanie do wskaźnika wiedzy', 'opisz, czego szukasz');
  const zakres = wybor('Zakres wskaźnika wiedzy', POZYCJE_ZAKRESU_WIEDZY);
  zakres.value = KnowledgeScope.All;
  const przebuduj = przelacznik('Przebuduj wskaźnik od zera');

  const lista = wykaz('Fragmenty odnalezione po znaczeniu', 'ma-wykaz');
  okno.tresc.append(lista);
  okno.puste(PUSTE.wiedzaSpoczynek);

  const szukaj = przyciskAkcji('Szukaj po znaczeniu', 'dn-btn dn-btn--sm dn-btn--atrament');
  szukaj.addEventListener('click', () => void szukajFragmentow());

  const zbuduj = przyciskAkcji('Zbuduj wskaźnik znaczenia', 'dn-btn dn-btn--sm dn-btn--atrament');
  zbuduj.addEventListener('click', () => void zbudujWskaznik());

  const dokumenty = przyciskAkcji('Wgraj dokument do bazy wiedzy', 'dn-btn dn-btn--sm dn-btn--duch');
  dokumenty.dataset['brak'] = 'wgranie-dokumentu';
  dokumenty.addEventListener('click', () => zglosBrak('Wgranie dokumentu', BRAKI.wgranieDokumentu));

  const szukanie = document.createElement('div');
  szukanie.className = 'ma-obszar';
  szukanie.dataset['obszar'] = 'pamiec-semantyczna';
  szukanie.append(
    wiersz('Pytanie', zapytanie, {
      klasa: 'ma-wiersz',
      objasnienie:
        'Szukanie idzie po sensie, nie po słowach (knowledge.search). Odpowiedzią są ' +
        'fragmenty wraz ze źródłem i trafnością, gotowe do zacytowania.',
    }),
    wiersz('Zakres', zakres, {
      klasa: 'ma-wiersz',
      objasnienie: 'Ten sam wskaźnik obsługuje bibliotekę, historię rozmów i pliki przestrzeni.',
    }),
    rzad(szukaj),
    okno.element,
    odpowiedz.element,
  );

  const wskaznik = document.createElement('div');
  wskaznik.className = 'ma-obszar';
  wskaznik.dataset['obszar'] = 'baza-wiedzy';
  wskaznik.append(
    wiersz('Przebuduj od zera', przebuduj, {
      klasa: 'ma-wiersz',
      objasnienie:
        'Bez przebudowy wskaźnik jest uzupełniany o pozycje nowe. Przebudowa liczy go ' +
        'od nowa i trwa dłużej.',
    }),
    rzad(zbuduj, dokumenty),
    opisBazy(),
  );

  async function szukajFragmentow(): Promise<void> {
    const pytanie = zapytanie.value.trim();
    if (pytanie === '') {
      odpowiedz.pokaz('Wpisz pytanie — rdzeń nie ma czego szukać bez frazy.', false);
      return;
    }
    okno.ladowanie(ODCZYTY.wiedza);
    const wynik = await zrodlo.szukaj({
      query: pytanie,
      scope: zakres.value as KnowledgeScope,
      ...(stan.idOkna() === '' ? {} : { windowId: stan.idOkna() }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Szukanie po znaczeniu', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const fragmenty = wynik.wynik.results;
    lista.replaceChildren(...fragmenty.map(pozycjaFragmentu));
    odpowiedz.pokaz(
      `Rdzeń zna ${String(wynik.wynik.total)} fragmentów spełniających warunki; ` +
        `pokazano ${String(fragmenty.length)}.`,
      true,
    );
    if (fragmenty.length === 0) {
      okno.puste(PUSTE.wiedza);
      return;
    }
    okno.gotowe();
  }

  async function zbudujWskaznik(): Promise<void> {
    okno.ladowanie(ODCZYTY.wskaznik);
    const wynik = await zrodlo.wskaznik({
      scope: zakres.value as KnowledgeScope,
      rebuild: przebuduj.checked,
      ...(stan.idOkna() === '' ? {} : { windowId: stan.idOkna() }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Budowa wskaźnika znaczenia', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const zbudowany = wynik.wynik;
    // Liczby pochodzą z odpowiedzi rdzenia, nie z zamówienia okna: wskaźnik
    // potrafi przyjąć mniej pozycji, niż zakres obejmuje.
    okno.puste(
      `Wskaźnik przyjął ${String(zbudowany.indexed)} pozycji i liczy ich teraz ` +
        `${String(zbudowany.total)}` +
        (zbudowany.model !== undefined && zbudowany.model !== ''
          ? ` (model osadzeń: ${zbudowany.model}).`
          : '. Rdzeń nie podał modelu osadzeń.') +
        ' Zapytaj w zakładce pamięci semantycznej, aby sięgnąć do wskaźnika.',
    );
  }

  return { szukanie, wskaznik };
}

/** Jeden odnaleziony fragment: treść, źródło i trafność podane przez rdzeń. */
function pozycjaFragmentu(trafienie: KnowledgeHit): HTMLElement {
  const czesci = [`źródło: ${trafienie.source}`];
  if (trafienie.sourceId !== undefined && trafienie.sourceId !== '') {
    czesci.push(`identyfikator źródła: ${trafienie.sourceId}`);
  }
  if (trafienie.score !== undefined) {
    czesci.push(`trafność: ${String(trafienie.score)}/100`);
  }
  czesci.push(`zakres: ${trafienie.scope}`);
  const { element } = pozycjaWykazu(trafienie.text, czesci.join(' · '), 'ma');
  return element;
}

/** Zdanie o tym, czym baza wiedzy profilu jest i skąd bierze treść. */
function opisBazy(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent =
    'Baza wiedzy profilu to wskaźnik znaczenia zbudowany z treści, które Operator ma już ' +
    'w platformie: biblioteki, historii rozmów i plików przestrzeni roboczej. Odpowiedzi ' +
    'asystenta cytują odnalezione fragmenty wraz z ich źródłem.';
  return element;
}

/** Rząd przycisków obszaru — ten sam odstęp w obu zakładkach. */
function rzad(...kontrolki: readonly HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'ma-formularz__przyciski';
  element.append(...kontrolki);
  return element;
}
