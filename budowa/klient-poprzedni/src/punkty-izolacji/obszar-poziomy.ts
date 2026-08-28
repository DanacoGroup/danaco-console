import { Command, ConfigScope, type IsolationScopeLevel, type RequestOf, type ResponseOf } from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import { utworzStanTresci } from '../komponenty/stan-tresci';
import { pole, wiersz as wierszPola } from '../modele/kontrolki-formularza-braki';
import type { Kanal, Wynik } from '../protokol/kanal';
import type { ObszarIzolacji, ZaleznosciObszaru } from './obszary';

// Obszar Poziom zasięgu — selektor zasięgu okna, drabina poziomów zapisu punktu izolacji.

/** Zdania objaśniające poziom — treść, której rdzeń dziś nie niesie, bo nie ma jej wcale u źródła danych. */
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
      // Zakres aplikacja obsługuje przełącznik logowania, stoi przy globalnym, bo dotyczy całej aplikacji.
  [ConfigScope.Application]: 'aplikacja',
  [ConfigScope.Global]:
    'Obejmuje całą platformę — warstwa bazowa, dziedziczona przez każdy węższy poziom ' +
    'bez własnego zapisu.',
};

const POWOD_BRAKU_ROZSTRZYGNIECIA =
  'isolation.scope.list zwraca słownik ośmiu poziomów, nie zapis Operatora na żadnym z nich. ' +
  'Rozstrzygnięcie „który poziom dziś wygrywa" wymaga pełnej ścieżki bytów (identyfikatora roli, ' +
  'sesji, projektu, pary modułów), której ten klient nie zna — tak samo jak w oknie Konfiguracji.';

/** Opakowuje wysłanie komendy kanału w obietnicę — kanał sam daje wyłącznie wersję z wywołaniem zwrotnym. */
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

  // Byt karty sesji klient zna sam; podpowiedź wchodzi wyłącznie do pola pustego, wpis własny nie ginie.
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
  // Kolejność rozstrzygania od najwęższego: pole kontraktu rośnie od najszerszego, wyświetlamy malejąco.
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

/** Notatka o braku rozstrzygnięcia: plakietka ostrzegawcza i zdanie z powodem, stojące zawsze nazwane wprost. */
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
