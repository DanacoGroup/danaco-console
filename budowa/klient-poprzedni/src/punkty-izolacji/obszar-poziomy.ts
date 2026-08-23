import { Command, ConfigScope, type IsolationScopeLevel, type RequestOf, type ResponseOf } from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import { utworzStanTresci } from '../komponenty/stan-tresci';
import { pole, wiersz as wierszPola } from '../modele/kontrolki-formularza-braki';
import type { Kanal, Wynik } from '../protokol/kanal';
import type { ObszarIzolacji, ZaleznosciObszaru } from './obszary';

/**
 * Obszar „Poziom zasięgu" — selektor zasięgu okna, czyli lewa kolumna rozdz. 6.1
 * Modelu konfiguracji: drabina poziomów zapisu punktu izolacji (okno → rola →
 * sesja → projekt → para modułów → moduł → środowisko → globalny), w kolejności
 * rozstrzygania — wygrywa pierwszy poziom mający własny zapis, zanim w grę
 * wchodzi oś (obszar `osie`).
 *
 * Drabina nie jest wykazem do czytania: wskazanie poziomu ustawia zasięg czynny
 * okna (`stan-zasiegu.ts`), a na nim czyta i pisze macierz izolacji, przypisanie
 * profilu i podgląd polityki efektywnej. Ograniczenie konfiguracji do jednego,
 * z góry narzuconego poziomu byłoby twardą regułą — a tych okno nie stawia.
 *
 * Drabina woła rdzeń (`isolation.scope.list`) i nie jest drugą kopią słownika.
 * Nazwa, kolejność i flaga najwęższego poziomu przychodzą z odpowiedzi rdzenia
 * (`shared.IsolationScopeLevel`), który czyta je z tabeli `poziom_zasiegu`
 * (`server/internal/core/adapter_modul_isolation.go`, funkcja `PoziomyZasiegu`)
 * — ten plik ich nie zgaduje ani nie trzyma jako stałej.
 *
 * Pole `order` w kontrakcie liczy odwrotnie niż kolejność rozstrzygania: rośnie
 * od poziomu najszerszego (`order: 1` dla globalnego), a tabela pokazuje
 * kolejność od najwęższego — stąd sortowanie malejąco po `order`.
 *
 * Czego rdzeń nie niesie: pola `description`, bo adapter pomija je tam, gdzie
 * baza nie ma treści. Zdania objaśniające przy każdym poziomie zostają więc
 * lokalnym słownikiem tego pliku (`ZDANIA`, kluczowany po `ConfigScope`
 * z kontraktu) — to nie jest druga kopia drabiny, tylko treść, której rdzeń
 * nie ma i nie udaje, że ma.
 *
 * Obszar nie pokazuje, który poziom wygrywa dla bieżącego Operatora — pokazuje
 * wyłącznie ten, który Operator wskazał do zapisu.
 * `isolation.scope.list` zwraca słownik poziomów (nazwa, kolejność, czy poziom
 * jest najwęższy w ogóle), a nie zapis Operatora na którymkolwiek z nich.
 * Takie rozstrzygnięcie wymaga innej rodziny komend (`isolation.context.get`,
 * `isolation.policy.preview`) oraz pełnej ścieżki bytów — identyfikatorów roli,
 * sesji, projektu i pary modułów — których klient nie zna. Bliźniacze okno
 * Konfiguracji ma tę samą drabinę (`konfiguracja/zasiegi.ts`) i z tego samego
 * powodu jego odczyt (`konfiguracja/zrodlo-wartosci.ts`, funkcja `wpisy`)
 * pobiera tylko dwa poziomy: globalny oraz ten wskazany punktem widzenia okna.
 * Braku nie wolno zamalować milczącą pustką, która wyglądałaby jak „nikt nic
 * nie zapisał" zamiast uczciwym „tego nie odczytaliśmy".
 */

/** Zdania objaśniające poziom — treść, której rdzeń dziś nie niesie (pole `description` puste u źródła). */
const ZDANIA: Readonly<Record<ConfigScope, string>> = {
  [ConfigScope.Window]:
    'Obejmuje wyłącznie jedno okno komunikacji — pojedynczą kartę rozmowy w module. ' +
    'Najwęższy poziom: mając własny zapis, wygrywa z każdym innym.',
  [ConfigScope.Role]:
    'Obejmuje jedną rolę przypisaną w panelu orkiestracji MultitaskingAI — np. profil ' +
    '„Executor 1" — a przez nią wszystkie okna tej roli.',
  [ConfigScope.Session]:
    'Obejmuje jedną kartę sesji — wszystkie okna otwarte w tej karcie; zapis wygasa ' +
    'z jej zamknięciem.',
  [ConfigScope.Project]: 'Obejmuje jeden projekt — np. rozdzielenie dwóch projektów otwartych w module Workspace.',
  [ConfigScope.ModulePair]:
    'Obejmuje wskazaną parę modułów — np. współdzielenie operacji AI między modułem ' +
    'Studio a modułem Translate.',
  [ConfigScope.Module]: 'Obejmuje jeden moduł platformy — np. odrębna pamięć długoterminowa modułu Developer.',
  [ConfigScope.Environment]: 'Obejmuje jedno środowisko uruchomieniowe — np. inna polityka dla CodeStudio niż dla TalkIn.',
      // Zakres „aplikacja" obsługuje przełącznik „Wymóg logowania". Stoi przy
      // globalnym, bo dotyczy całej aplikacji — szerzej się nie da.
  [ConfigScope.Application]: 'aplikacja',
  [ConfigScope.Global]:
    'Obejmuje całą platformę — warstwa bazowa, dziedziczona przez każdy węższy poziom ' +
    'bez własnego zapisu.',
};

const POWOD_BRAKU_ROZSTRZYGNIECIA =
  'isolation.scope.list zwraca słownik ośmiu poziomów, nie zapis Operatora na żadnym z nich. ' +
  'Rozstrzygnięcie „który poziom dziś wygrywa" wymaga pełnej ścieżki bytów (identyfikatora roli, ' +
  'sesji, projektu, pary modułów), której ten klient nie zna — tak samo jak w oknie Konfiguracji.';

/** Opakowuje `kanal.wyslij` w Promise — kanał sam daje wyłącznie wersję z wywołaniem zwrotnym. */
function posijKomende<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    kanal.wyslij(komenda, zadanie, (wynik) => rozstrzygnij(wynik));
  });
}

export function utworzObszar(zaleznosci: ZaleznosciObszaru): ObszarIzolacji {
  const { kanal, zasieg } = zaleznosci;
  const stan = utworzStanTresci('pi');

  const bytPoziomu = pole('Identyfikator bytu poziomu', 'puste dla poziomu globalnego');
  bytPoziomu.value = zasieg.bytZasiegu();
  bytPoziomu.addEventListener('change', () => zasieg.ustaw(zasieg.zasieg(), bytPoziomu.value.trim()));

  /**
   * Byt karty sesji klient zna sam — wpisywanie go ręcznie z pamięci byłoby
   * proszeniem Operatora o identyfikator, który okno ma pod ręką. Podpowiedź
   * wchodzi wyłącznie do pola pustego: wpis własny nie jest nadpisywany.
   */
  function wskazPoziom(wybrany: ConfigScope): void {
    const wpisany = bytPoziomu.value.trim();
    const byt = wpisany === '' && wybrany === ConfigScope.Session ? kanal.sesja().id() : wpisany;
    bytPoziomu.value = byt;
    zasieg.ustaw(wybrany, byt);
  }

  async function odswiez(): Promise<void> {
    stan.ladowanie('Pytam rdzeń o drabinę poziomów zasięgu (isolation.scope.list)…');

    const idSesji = kanal.sesja().id();
    const wynik = await posijKomende(kanal, Command.IsolationScopeList, {
      sessionId: idSesji || undefined,
    });

    if (!wynik.udany || wynik.wynik === undefined) {
      stan.blad(opisOdmowyBledu('Odczyt poziomów zasięgu', wynik.blad), wynik.blad);
      return;
    }

    if (wynik.wynik.scopes.length === 0) {
      stan.pusto('Rdzeń nie zwrócił żadnego poziomu zasięgu.');
      return;
    }

    const miejsce = stan.tresc();
    miejsce.append(
      zbudujWstep(),
      zbudujTabeleDrabiny(wynik.wynik.scopes, zasieg.zasieg(), wskazPoziom),
      wierszPola('Identyfikator bytu poziomu', bytPoziomu, { klasa: 'dn-pole' }),
      zbudujNotatkeBrakuRozstrzygniecia(),
    );
    stan.potwierdzenie(
      `Zasięg czynny okna: ${zasieg.opis()}. Na nim czyta i pisze macierz izolacji.`,
      true,
    );
  }

  return {
    element: stan.element,
    odswiez: () => void odswiez(),
  };
}

function zbudujWstep(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent =
    'Wskaż poziom, na którym reguła izolacji ma obowiązywać — macierz izolacji, przypisanie profilu ' +
    'i podgląd polityki czytają i piszą właśnie na nim. Poziom zasięgu rozstrzyga PIERWSZY, zanim ' +
    'w grę wejdzie oś rozstrzygania: rdzeń sprawdza poziomy w tej kolejności i zatrzymuje się na ' +
    'pierwszym, który ma własny zapis.';
  return element;
}

function zbudujTabeleDrabiny(
  scopes: readonly IsolationScopeLevel[],
  wybrany: ConfigScope,
  naWybor: (poziom: ConfigScope) => void,
): HTMLElement {
  // Kolejność rozstrzygania od najwęższego: `order` kontraktu rośnie od
  // najszerszego (1 = globalny), więc wyświetlamy malejąco po `order`.
  const uporzadkowane = [...scopes].sort((a, b) => b.order - a.order);

  const naglowek = document.createElement('tr');
  naglowek.append(
    komorkaNaglowka('Wybór'),
    komorkaNaglowka('Kolejność'),
    komorkaNaglowka('Poziom'),
    komorkaNaglowka('Czego dotyczy zapis'),
  );

  const glowa = document.createElement('thead');
  glowa.append(naglowek);

  const cialo = document.createElement('tbody');
  cialo.append(...uporzadkowane.map((wiersz, indeks) => zbudujWiersz(wiersz, indeks + 1, wybrany, naWybor)));

  const tabela = document.createElement('table');
  tabela.className = 'dn-tabela';
  tabela.setAttribute('aria-label', 'Osiem poziomów zasięgu w kolejności rozstrzygania, od najwęższego');
  tabela.append(glowa, cialo);
  return tabela;
}

function komorkaNaglowka(tytul: string): HTMLElement {
  const komorka = document.createElement('th');
  komorka.scope = 'col';
  komorka.textContent = tytul;
  return komorka;
}

function zbudujWiersz(
  wiersz: IsolationScopeLevel,
  kolejnosc: number,
  wybrany: ConfigScope,
  naWybor: (poziom: ConfigScope) => void,
): HTMLElement {
  const element = document.createElement('tr');
  element.dataset['poziom'] = wiersz.scope;
  if (wiersz.narrowest === true) element.dataset['najwezszy'] = 'true';
  if (wiersz.scope === wybrany) element.dataset['wybrany'] = 'true';

  const wybor = document.createElement('input');
  wybor.type = 'radio';
  wybor.name = 'pi-poziom-zasiegu';
  wybor.value = wiersz.scope;
  wybor.checked = wiersz.scope === wybrany;
  wybor.setAttribute('aria-label', `Ustaw zasięg na poziom „${wiersz.label}"`);
  wybor.addEventListener('change', () => naWybor(wiersz.scope));

  const wyborKomorka = document.createElement('td');
  wyborKomorka.append(wybor);

  const kolejnoscKomorka = document.createElement('td');
  kolejnoscKomorka.className = 'dn-dane';
  kolejnoscKomorka.textContent = String(kolejnosc);

  const nazwaKomorka = document.createElement('td');
  nazwaKomorka.textContent = wiersz.label;

  const zdanieKomorka = document.createElement('td');
  zdanieKomorka.className = 'dn-dane';
  zdanieKomorka.textContent =
    wiersz.description ?? ZDANIA[wiersz.scope] ?? 'Rdzeń nie podał objaśnienia tego poziomu.';

  element.append(wyborKomorka, kolejnoscKomorka, nazwaKomorka, zdanieKomorka);
  return element;
}

/**
 * Notatka o braku rozstrzygnięcia: plakietka ostrzegawcza i zdanie z powodem —
 * `isolation.scope.list` niesie słownik poziomów, nie zapis Operatora na
 * żadnym z nich, więc obszar nie umie wskazać, który poziom wygrywa dla
 * bieżącego okna albo sesji. Braku nie chowamy za `disabled` ani za milczącą
 * pustką — stoi nazwany wprost.
 */
function zbudujNotatkeBrakuRozstrzygniecia(): HTMLElement {
  const plakietka = document.createElement('span');
  plakietka.className = 'dn-plakietka dn-plakietka--ostrzezenie';
  plakietka.textContent = 'Rozstrzygnięcie dla bieżącego okna niedostępne';

  const powod = document.createElement('p');
  powod.className = 'dn-pole-opis';
  powod.textContent = POWOD_BRAKU_ROZSTRZYGNIECIA;

  const element = document.createElement('div');
  element.append(plakietka, powod);
  return element;
}
