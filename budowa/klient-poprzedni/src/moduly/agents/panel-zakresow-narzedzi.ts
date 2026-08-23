import type { ToolScope } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { ZrodloZakresuEksperta } from './zrodlo-zakresu-eksperta';

/**
 * Panel zakresów narzędzi profilu asystenta (`tools.scope.*`).
 *
 * Stoi w Permissions Center, choć dotyczy profilu asystenta, a nie eksperta.
 * Powód jest jeden i praktyczny: to jedyne okno, w którym Operator ustala „jak
 * szeroko wykonawca działa w moim imieniu", więc rozstrzygnięcie o wywoływaniu
 * narzędzi ma stać tam, gdzie Operator go szuka.
 *
 * ZAKRES JEST NASTAWĄ ZASIĘGU, NIE BRAMKĄ WBUDOWANĄ. Pozycja bez wiersza jest
 * dostępna i nie ma granicy wywołań — platforma niczego nie zawęża z góry.
 * Zawężenie zapisane tutaj rdzeń czyta przed każdym wywołaniem ręki modelu
 * i odmawia; kolumna zużycia pokazuje, ile z granicy zostało.
 */
export interface PanelZakresowNarzedzi {
  element: HTMLElement;
  /** Odczytuje zakresy profilu wskazanego w polu; puste bierze profil domyślny. */
  wczytaj(): Promise<void>;
}

export function utworzPanelZakresowNarzedzi(
  zrodlo: ZrodloZakresuEksperta,
): PanelZakresowNarzedzi {
  const odpowiedz = utworzWierszOdpowiedzi();

  const profil = poleTekstowe({
    etykieta: 'Profil asystenta',
    podpowiedz: 'puste bierze profil domyślny',
    opis: 'Zakresy zapisują się przy profilu, nie przy oknie: profil jest tym, co Operator wybiera dla asystenta.',
  });

  const pozycja = poleTekstowe({
    etykieta: 'Pozycja katalogu — nazwa pełna ze źródłem',
    podpowiedz: 'np. danaco:session.rename',
    opis: 'Nazwa pełna, bo skrócona nie jest jednoznaczna — ta sama może przyjść z dwóch źródeł.',
  });

  const dostepna = document.createElement('input');
  dostepna.type = 'checkbox';
  dostepna.className = 'dn-przelacznik';
  dostepna.checked = true;

  const potwierdzenie = document.createElement('input');
  potwierdzenie.type = 'checkbox';
  potwierdzenie.className = 'dn-przelacznik';

  const limit = document.createElement('input');
  limit.type = 'number';
  limit.className = 'dn-pole-kontrolka';
  limit.min = '0';
  limit.value = '0';

  const okno = document.createElement('input');
  okno.type = 'number';
  okno.className = 'dn-pole-kontrolka';
  okno.min = '1';
  okno.value = '3600';

  const zapisz = przycisk('Zapisz zakres pozycji', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odczytaj = przycisk('Odczytaj zakresy profilu', 'dn-btn dn-btn--sm dn-btn--zarys');

  const nastawy = document.createElement('div');
  nastawy.className = 'dn-pole dm-pole';
  const opisNastaw = document.createElement('p');
  opisNastaw.className = 'dn-pole-opis';
  opisNastaw.textContent =
    'Dostępna · wymaga potwierdzenia · granica wywołań (zero znaczy bez granicy) · okno czasu ' +
    'w sekundach. Granica jest egzekwowana: po jej wyczerpaniu rdzeń odmawia wywołania.';
  nastawy.append(dostepna, potwierdzenie, limit, okno, opisNastaw);

  const wykaz = document.createElement('ul');
  wykaz.className = 'da-uprawnienia da-uprawnienia--narzedzia';

  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Zakres narzędzi profilu asystenta';

  const element = document.createElement('section');
  element.className = 'da-panel da-panel--narzedzia';
  element.dataset['panel'] = 'zakresy-narzedzi';
  element.append(
    tytul,
    profil.element,
    pozycja.element,
    nastawy,
    zapisz,
    odczytaj,
    wykaz,
    odpowiedz.element,
  );

  function wierszZakresu(zakres: ToolScope): HTMLElement {
    const nazwa = document.createElement('span');
    nazwa.className = 'da-uprawnienia__etykieta';
    nazwa.textContent = `${zakres.shortName} (${zakres.toolName})`;

    const stan = document.createElement('span');
    stan.className = 'da-uprawnienia__zakresy';
    const zuzycie =
      zakres.callLimit === 0
        ? 'bez granicy'
        : `${zakres.callsUsed ?? 0} z ${zakres.callLimit} w oknie ${zakres.callWindowSeconds} s`;
    stan.textContent =
      `${zakres.enabled ? 'dostępna' : 'WYŁĄCZONA'} · ` +
      `${zakres.confirmRequired ? 'wymaga potwierdzenia' : 'bez potwierdzenia'} · ${zuzycie}`;

    const wiersz = document.createElement('li');
    wiersz.className = 'da-uprawnienia__wiersz';
    wiersz.dataset['pozycja'] = zakres.toolName;
    wiersz.append(nazwa, stan);
    return wiersz;
  }

  async function wczytaj(): Promise<void> {
    const wynik = await zrodlo.zakresyNarzedzi(profil.kontrolka.value.trim());
    if (!wynik.udany || wynik.wynik === undefined) {
      wykaz.replaceChildren();
      odpowiedz.pokaz(
        opisOdmowy('Odczyt zakresów narzędzi', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    wykaz.replaceChildren(...wynik.wynik.scopes.map(wierszZakresu));
    odpowiedz.pokaz(
      wynik.wynik.total === 0
        ? 'Profil nie ma zapisanego ani jednego zawężenia — każda pozycja katalogu jest dostępna bez granicy.'
        : `Profil ma ${wynik.wynik.total} zapisanych zakresów.`,
      true,
    );
  }

  async function zapiszZakres(): Promise<void> {
    const nazwa = pozycja.kontrolka.value.trim();
    if (nazwa === '') {
      odpowiedz.pokaz('Podaj nazwę pełną pozycji katalogu — zakres bez pozycji nie ma o czym mówić.', false);
      return;
    }
    const idProfilu = profil.kontrolka.value.trim();
    if (idProfilu === '') {
      odpowiedz.pokaz(
        'Zapis zakresu wymaga wskazania profilu wprost — odczyt bierze domyślny, zapis nie zgaduje.',
        false,
      );
      return;
    }
    odpowiedz.pokaz(`Zapis zakresu pozycji ${nazwa}…`, true);
    const wynik = await zrodlo.zapiszZakresNarzedzia({
      idProfilu,
      nazwaPozycji: nazwa,
      dostepna: dostepna.checked,
      potwierdzenie: potwierdzenie.checked,
      granicaWywolan: Number.parseInt(limit.value, 10) || 0,
      oknoSekund: Number.parseInt(okno.value, 10) || 3600,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Zapis zakresu narzędzia', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    await wczytaj();
    const zapisany = wynik.wynik.scope;
    odpowiedz.pokaz(
      zapisany.enabled
        ? `Pozycja ${zapisany.shortName} dostępna profilowi` +
            (zapisany.callLimit === 0
              ? ' bez granicy wywołań.'
              : `, granica ${zapisany.callLimit} wywołań w oknie ${zapisany.callWindowSeconds} s.`)
        : `Pozycja ${zapisany.shortName} WYŁĄCZONA dla tego profilu — rdzeń odmówi jej wywołania.`,
      true,
    );
  }

  zapisz.addEventListener('click', () => void zapiszZakres());
  odczytaj.addEventListener('click', () => void wczytaj());

  return { element, wczytaj };
}
