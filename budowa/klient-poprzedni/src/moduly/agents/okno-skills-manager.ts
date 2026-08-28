import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { Kanal } from '../../protokol/kanal';
import { utworzZrodloZakresuEksperta, type ZrodloZakresuEksperta } from './zrodlo-zakresu-eksperta';
import { naglowek } from './biblioteka-ekspertow';
import { utworzSterNarzedzi, type SterNarzedzi } from './drzewo-narzedzi';
import { GRUPY_NARZEDZI, KATALOG_NARZEDZI } from './katalog-narzedzi';
import { utworzLicznikNarzedzi, type LicznikNarzedzi } from './licznik-narzedzi';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAgentow } from './stan-agentow';

/**
 * Skills Manager jest oknem zarządcą modułu Agents: dobiera narzędzia
 * z wykazu pogrupowanego po przeznaczeniu, wraz z polem otwartym na kody
 * spoza katalogu kontraktu.
 */
export interface OknoSkillsManager {
  element: HTMLElement;
  odswiez(): void;
  /** Zwija drzewo wyboru i zdejmuje jego nasłuchy dokumentu. */
  zamknij(): void;
}

export function utworzOknoSkillsManager(
  stan: StanAgentow,
  kanal: Kanal,
): OknoSkillsManager {
  const okno: StanOkna = utworzStanOkna();
  const zakres: ZrodloZakresuEksperta = utworzZrodloZakresuEksperta(kanal);

  const lista = document.createElement('ul');
  lista.className = 'da-umiejetnosci';

  const licznik: LicznikNarzedzi = utworzLicznikNarzedzi();
  const ster: SterNarzedzi = utworzSterNarzedzi();

  const podpowiedzi = document.createElement('datalist');
  podpowiedzi.id = 'da-umiejetnosci-podpowiedzi';

  const pole = poleTekstowe({
    etykieta: 'Kod do przypisania',
    podpowiedz: 'nazwa narzędzia albo obszaru',
    opis:
      'Pole i drzewo pokazują tę samą nastawę. Pole zostaje otwarte, bo rdzeń ' +
      'przyjmuje też kody spoza katalogu kontraktu — wpisz taki wprost.',
  });
  pole.kontrolka.setAttribute('list', podpowiedzi.id);

  const dodaj = przycisk('Przypisz kod ekspertowi', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odpowiedz = utworzWierszOdpowiedzi();

  const wybor = document.createElement('div');
  wybor.className = 'da-wybor-narzedzi';
  wybor.dataset['wybor'] = 'narzedzia';

  const zdanieWyboru = document.createElement('p');
  zdanieWyboru.className = 'dn-pole-opis';
  zdanieWyboru.textContent =
    `Wykaz kontraktu: ${KATALOG_NARZEDZI.length} narzędzi w ${GRUPY_NARZEDZI.length} ` +
    'obszarach. Gałąź jest obszarem, a jej pierwsza pozycja przypisuje cały obszar ' +
    'jednym kodem — tak samo czyta te kody serwer narzędzi.';

  wybor.append(zdanieWyboru, ster.element);

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis da-granica';
  granica.textContent =
    'Granica okna: nazwa gałęzi jest kodem obszaru, nie zdaniem „do czego służy" — ' +
    'opisy obszarów stoją w źródle kontraktu, lecz generator nie oddaje ich klientowi. ' +
    'Odłączenie kodu ma własną kontrolkę przy każdym wierszu wykazu wyżej i to ona ' +
    'nazywa swój brak.';

  okno.tresc.append(lista, wybor, pole.element, podpowiedzi, dodaj, odpowiedz.element, granica);

  const element = document.createElement('section');
  element.className = 'da-okno da-okno--zarzadca';
  element.dataset['okno'] = 'skills-manager';
  element.append(naglowek('Skills Manager'), licznik.element, okno.element);

  // Dwie kontrolki, jeden kod: żadna ze stron nie ogłasza zmiany drugiej z powrotem, więc pętli nie ma.
  ster.naZmiane((kod) => {
    pole.kontrolka.value = kod;
  });
  pole.kontrolka.addEventListener('input', () => {
    ster.ustawKod(pole.kontrolka.value.trim());
  });

  async function przypisz(): Promise<void> {
    const ekspert = stan.wybrany();
    if (ekspert === null) {
      odpowiedz.pokaz('Wybierz eksperta w Agent Builderze — kod należy do jednego.', false);
      return;
    }
    const umiejetnosc = pole.kontrolka.value.trim();
    if (umiejetnosc === '') {
      odpowiedz.pokaz(
        'Wskaż kod w drzewie albo wpisz go w polu — rdzeń odmówi przypisania bez niego.',
        false,
      );
      return;
    }
    const mialWczesniej = (ekspert.skillIds ?? []).includes(umiejetnosc);
    odpowiedz.pokaz(`Przypisywanie kodu ${umiejetnosc}…`, true);
    const wynik = await stan.zrodlo.dodajUmiejetnosc(ekspert.id, umiejetnosc);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Przypisanie kodu', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const zapisany = wynik.wynik.agent;
    stan.wchlon(zapisany);
    pole.kontrolka.value = '';
    ster.ustawKod('');
    const maTeraz = (zapisany.skillIds ?? []).includes(umiejetnosc);
    if (mialWczesniej || !maTeraz) {
      odpowiedz.pokaz(
        mialWczesniej
          ? `Ekspert „${zapisany.name}" miał już kod ${umiejetnosc} — nic się nie zmieniło.`
          : `Rdzeń przyjął wywołanie, ale kodu ${umiejetnosc} nie ma w wykazie ` +
              `eksperta „${zapisany.name}" — przypisanie się nie odbyło.`,
        false,
      );
      return;
    }
    odpowiedz.pokaz(`Kod ${umiejetnosc} przypisany ekspertowi „${zapisany.name}".`, true);
  }

  /** Zdjęcie umiejętności z definicji eksperta; świadectwem powodzenia jest wykaz, który wrócił. */
  async function odlacz(umiejetnosc: string): Promise<void> {
    const ekspert = stan.wybrany();
    if (ekspert === null) return;
    odpowiedz.pokaz(`Zdejmowanie kodu ${umiejetnosc}…`, true);
    const wynik = await zakres.odlaczUmiejetnosc(ekspert.id, umiejetnosc);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Zdjęcie kodu', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const zapisany = wynik.wynik.agent;
    stan.wchlon(zapisany);
    if ((zapisany.skillIds ?? []).includes(umiejetnosc)) {
      odpowiedz.pokaz(
        `Rdzeń przyjął wywołanie, ale kod ${umiejetnosc} nadal stoi w wykazie ` +
          `eksperta „${zapisany.name}" — zdjęcie się nie odbyło.`,
        false,
      );
      return;
    }
    odpowiedz.pokaz(`Kod ${umiejetnosc} zdjęty z definicji eksperta „${zapisany.name}".`, true);
  }

  dodaj.addEventListener('click', () => void przypisz());

  /** Ekspert, którego wykaz stoi w oknie — po nim poznajemy zmianę wyboru. */
  let pokazany = '';

  function odswiez(): void {
    const ekspert = stan.wybrany();
    licznik.ustaw(ekspert);
    ustawPodpowiedzi();
    // Kod wpisany dla jednego eksperta nie może czekać w polu przy drugim ekspercie.
    if ((ekspert?.id ?? '') !== pokazany) {
      pokazany = ekspert?.id ?? '';
      pole.kontrolka.value = '';
      ster.ustawKod('');
      odpowiedz.wyczysc();
    }
    // Drzewo znakuje to, co ekspert już ma, obiema listami kodów naraz.
    ster.ustawPrzypisane([...(ekspert?.skillIds ?? []), ...(ekspert?.connectorIds ?? [])]);
    if (ekspert === null) {
      lista.replaceChildren();
      okno.puste('Wybierz eksperta w Agent Builderze, aby zobaczyć jego wykaz kodów.');
      return;
    }
    const umiejetnosci = ekspert.skillIds ?? [];
    lista.replaceChildren(
      ...umiejetnosci.map((kod) => wiersz(kod, (wskazany) => void odlacz(wskazany))),
    );
    if (umiejetnosci.length === 0) {
      okno.puste(
        `Ekspert „${ekspert.name}" nie ma jeszcze przypisanego ani jednego kodu — ` +
          'dopóki go nie dostanie, model widzi wykaz narzędzi w całości.',
      );
      return;
    }
    okno.gotowe();
  }

  /** Podpowiedź pola otwartego zostaje przy tym, czego drzewo nie zna — kodach spoza katalogu kontraktu. */
  function ustawPodpowiedzi(): void {
    const zKatalogu = new Set(KATALOG_NARZEDZI.map((pozycja) => pozycja.nazwa));
    const znane = new Set<string>();
    for (const ekspert of stan.eksperci()) {
      for (const kod of [...(ekspert.skillIds ?? []), ...(ekspert.connectorIds ?? [])]) {
        if (kod !== '' && !zKatalogu.has(kod) && !GRUPY_NARZEDZI.includes(kod)) znane.add(kod);
      }
    }
    podpowiedzi.replaceChildren(
      ...[...znane].sort((pierwsza, druga) => pierwsza.localeCompare(druga, 'pl')).map((kod) => {
        const pozycja = document.createElement('option');
        pozycja.value = kod;
        return pozycja;
      }),
    );
  }



  return {
    element,
    odswiez,
    zamknij: () => {
      ster.zwin();

    },
  };
}

/**
 * Wiersz przypisanego kodu wraz z kontrolką odłączenia, widoczną i klikalną
 * niezależnie od tego, czy odłączenie ma dziś drogę do rdzenia.
 */
function wiersz(kod: string, odlacz: (kod: string) => void): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'da-umiejetnosci__kod';
  nazwa.textContent = kod;

  const element = document.createElement('li');
  element.className = 'da-umiejetnosci__wiersz';
  element.dataset['umiejetnosc'] = kod;
  element.append(
    nazwa,
    przyciskOdlaczenia(kod, odlacz),
  );
  return element;
}

/** Kontrolka odłączenia umiejętności, osobna na każdy wiersz wykazu przypisanych kodów tego eksperta w bibliotece. */
function przyciskOdlaczenia(kod: string, odlacz: (kod: string) => void): HTMLButtonElement {
  const kontrolka = przycisk('Odłącz', 'dn-btn dn-btn--sm dn-btn--zarys');
  kontrolka.title = `zdjęcie kodu ${kod} z definicji tego eksperta`;
  kontrolka.addEventListener('click', () => odlacz(kod));
  return kontrolka;
}
