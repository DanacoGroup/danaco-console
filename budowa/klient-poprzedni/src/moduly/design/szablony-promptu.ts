import { Command, type DesignPrompt } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import type { StanDesignu } from './stan-designu';

/**
 * Szablony i historia promptów w Prompt Builderze, trwałe na koncie zamiast ginące z zamknięciem
 * karty przeglądarki.
 */
export interface SzablonyPromptu {
  element: HTMLElement;
  /** Zleca odczyt szablonów i historii okna. */
  wczytaj(): Promise<void>;
}

/** Czym panel szablonów i historii promptów steruje wprost w oknie kreatora, wraz z jego bieżącą treścią. */
export interface SterowanieSzablonami {
  /** Prompt zbudowany w kreatorze — zapisywany jako szablon. */
  prompt(): DesignPrompt;
  /** Wstawia prompt z szablonu albo z historii w pola kreatora. */
  naSzablon(prompt: DesignPrompt): void;
}

export function utworzSzablonyPromptu(
  stan: StanDesignu,
  sterowanie: SterowanieSzablonami,
): SzablonyPromptu {
  let szablony: readonly { id: string; name: string; prompt: DesignPrompt }[] = [];

  const nazwa = poleTekstowe({
    etykieta: 'Nazwa szablonu',
    podpowiedz: 'np. Bohater strony — wariant złoty',
    opis: 'Pole name żądania design.prompt.template.save.',
  });

  const wybor = poleWyboru(
    {
      etykieta: 'Szablon',
      opis:
        'Wybrany szablon zasila pola kreatora i staje się celem nadpisania ' +
        '(pole templateId). Pozycja „nowy szablon" zakłada kolejny obok.',
    },
    [{ wartosc: '', etykieta: 'nowy szablon' }],
  );

  const zapisz = przycisk('Zapisz jako szablon', 'dn-btn dn-btn--sm dn-btn--atrament');
  const wczytajSzablony = przycisk('Odczytaj szablony', 'dn-btn dn-btn--sm dn-btn--zarys');
  const uzyj = przycisk('Wstaw szablon w pola', 'dn-btn dn-btn--sm dn-btn--zarys');
  const wczytajHistorie = przycisk('Odczytaj historię poleceń', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const pasek = document.createElement('div');
  pasek.className = 'md-czynnosci__pasek';
  pasek.append(zapisz, wczytajSzablony, uzyj, wczytajHistorie);

  const historia = document.createElement('ul');
  historia.className = 'md-historia__wykaz';

  const element = document.createElement('div');
  element.className = 'md-szablony';
  element.append(nazwa.element, wybor.element, pasek, historia, odpowiedz.element);

  zapisz.addEventListener('click', () => void zapiszSzablon());
  wczytajSzablony.addEventListener('click', () => void odczytajSzablony());
  wczytajHistorie.addEventListener('click', () => void odczytajHistorie());
  uzyj.addEventListener('click', () => wstawSzablon());
  // Wybór pozycji przepisuje nazwę do pola: nadpisanie bez nazwy zapisałoby szablon pod nazwą pustą.
  wybor.kontrolka.addEventListener('change', () => {
    const szablon = szablony.find((pozycja) => pozycja.id === wybor.kontrolka.value);
    nazwa.kontrolka.value = szablon?.name ?? '';
  });

  function pokazSzablony(): void {
    ustawPozycje(wybor.kontrolka, [
      { wartosc: '', etykieta: 'nowy szablon' },
      ...szablony.map((pozycja) => ({ wartosc: pozycja.id, etykieta: pozycja.name })),
    ]);
  }

  function wstawSzablon(): void {
    const szablon = szablony.find((pozycja) => pozycja.id === wybor.kontrolka.value);
    if (szablon === undefined) {
      odpowiedz.pokaz('Wskaż szablon — pozycja „nowy szablon" nie ma czego wstawić.', false);
      return;
    }
    sterowanie.naSzablon(szablon.prompt);
    odpowiedz.pokaz(`Pola kreatora niosą teraz szablon „${szablon.name}".`, true);
  }

  async function zapiszSzablon(): Promise<void> {
    if (stan.idOkna() === '') {
      odpowiedz.pokaz(
        `Komenda ${Command.DesignPromptTemplateSave} wymaga okna modułu. ${stan.opisOkna()}`,
        false,
      );
      return;
    }
    const prompt = sterowanie.prompt();
    if (prompt.subject.trim() === '') {
      odpowiedz.pokaz(
        'Szablon bez tematu nie wygeneruje niczego, gdy Operator po niego sięgnie — ' +
          'wypełnij pole tematu w kreatorze.',
        false,
      );
      return;
    }
    odpowiedz.pokaz('Zapis szablonu…', true);
    const wynik = await stan.czuwanie.prowadz(
      'zapis szablonu promptu',
      stan.zrodlo.zapiszSzablonPromptu({
        idOkna: stan.idOkna(),
        nazwa: nazwa.kontrolka.value === '' ? prompt.subject : nazwa.kontrolka.value,
        prompt,
        idSzablonu: wybor.kontrolka.value,
      }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Zapis szablonu promptu', wynik.blad), false);
      return;
    }
    const szablon = wynik.wynik.template;
    szablony = [
      { id: szablon.id, name: szablon.name, prompt: szablon.prompt },
      ...szablony.filter((pozycja) => pozycja.id !== szablon.id),
    ];
    pokazSzablony();
    wybor.kontrolka.value = szablon.id;
    odpowiedz.pokaz(`Rdzeń utrwalił szablon „${szablon.name}" (${szablon.id}).`, true);
  }

  async function odczytajSzablony(): Promise<void> {
    if (stan.idOkna() === '') {
      odpowiedz.pokaz(
        `Komenda ${Command.DesignPromptTemplateList} wymaga okna modułu. ${stan.opisOkna()}`,
        false,
      );
      return;
    }
    odpowiedz.pokaz('Odczyt szablonów okna…', true);
    const wynik = await stan.czuwanie.prowadz(
      'odczyt szablonów promptów',
      stan.zrodlo.szablonyPromptu(stan.idOkna()),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt szablonów promptów', wynik.blad), false);
      return;
    }
    szablony = wynik.wynik.templates.map((szablon) => ({
      id: szablon.id,
      name: szablon.name,
      prompt: szablon.prompt,
    }));
    pokazSzablony();
    odpowiedz.pokaz(
      szablony.length === 0
        ? 'Rdzeń nie zna ani jednego szablonu tego okna.'
        : `Szablonów w oknie: ${wynik.wynik.total}.`,
      true,
    );
  }

  async function odczytajHistorie(): Promise<void> {
    if (stan.idOkna() === '') {
      odpowiedz.pokaz(
        `Komenda ${Command.DesignPromptHistoryList} wymaga okna modułu. ${stan.opisOkna()}`,
        false,
      );
      return;
    }
    odpowiedz.pokaz('Odczyt historii poleceń…', true);
    const wynik = await stan.czuwanie.prowadz(
      'odczyt historii promptów',
      // Granica zero zostawia rdzeniowi jego własną: przycięcie wymyślone tutaj ukrywałoby część historii.
      stan.zrodlo.historiaPromptow(stan.idOkna(), 0),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt historii promptów', wynik.blad), false);
      return;
    }
    historia.replaceChildren(
      ...wynik.wynik.prompts.map((zapis) => {
        const wiersz = document.createElement('li');
        wiersz.className = 'md-historia__pozycja';
        wiersz.dataset['prompt'] = zapis.id;
        const zasobow = zapis.assetIds?.length ?? 0;
        wiersz.textContent =
          `${zapis.prompt.subject} — zasobów: ${zasobow}` +
          (zapis.channelId === undefined ? '' : ` · kanał ${zapis.channelId}`);
        // Pozycja historii jest klikalna: wraca w pola kreatora, żeby polecenie dało się poprawić i wydać.
        wiersz.tabIndex = 0;
        wiersz.addEventListener('click', () => {
          sterowanie.naSzablon(zapis.prompt);
          odpowiedz.pokaz(`Pola kreatora niosą teraz polecenie z ${zapis.id}.`, true);
        });
        return wiersz;
      }),
    );
    odpowiedz.pokaz(
      wynik.wynik.total === 0
        ? 'W tym oknie nie wydano jeszcze ani jednego polecenia.'
        : `Poleceń wydanych w oknie: ${wynik.wynik.total}.`,
      true,
    );
  }

  async function wczytaj(): Promise<void> {
    await odczytajSzablony();
  }

  return { element, wczytaj };
}
