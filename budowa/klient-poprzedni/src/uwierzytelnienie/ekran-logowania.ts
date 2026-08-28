import './uwierzytelnienie.css';

import type { AuthSession } from '../../../shared/contract';
import { utworzStanTresci } from '../komponenty/stan-tresci';
import type { Kanal } from '../protokol/kanal';
import { opisOdmowyBramki } from './odmowy-auth';
import {
  etykietaPierwszego,
  MIN_ZNAKOW,
  NAPISY_PRZYCISKU,
  OBJASNIENIA,
  poleHasla,
  poleTekstu,
  utworzNiewylogowuj,
  utworzOdnosnikPotwierdzenia,
  utworzOdnosnikResetu,
  utworzSegmentyMetody,
  zlozPrzeslone,
  type Tryb,
} from './postac-bramki';
import { rozpoznajWejscie } from './rozpoznanie-bramki';
import {
  odczytajSesje,
  opisWaznosci,
  sesjaTrwala,
  skasujSesje,
  zapiszSesje,
} from './sesja-bramki';
import { odczytajTozsamoscUrzadzenia } from './tozsamosc-urzadzenia';
import { bezUchwytu, utworzZrodloUwierzytelnienia, type MetodaWejscia } from './zrodlo-auth';

/** Ekran logowania to jedyna bramka produktu — przesłona nad złożoną już aplikacją, nie osobna trasa; po wejściu znika, a rdzeń rozstrzyga wybór formularza i czas ważności sesji. */
export interface OpisEkranuLogowania {
  kanal: Kanal;
  /** Wejście udane — aplikacja pod przesłoną dostaje sesję bramki. */
  naWejscie(sesja: AuthSession): void;
}

export interface EkranLogowania {
  element: HTMLElement;
  /** Montuje przesłonę i rozpoczyna rozpoznanie: zapisana sesja → stan bramki → formularz. */
  uruchom(): void;
  rozlacz(): void;
}

export function utworzEkranLogowania(opis: OpisEkranuLogowania): EkranLogowania {
  const zrodlo = utworzZrodloUwierzytelnienia(opis.kanal);
  const stany = utworzStanTresci('au');

  let tryb: Tryb = 'wejscie';
  let metoda: MetodaWejscia = 'haslo';
  let metody: MetodaWejscia[] = ['haslo'];
  let wpuszczony = false;

  // ── formularz budowany raz; tryb przełącza etykiety i pole powtórzenia ─────

  const haslo = poleHasla('au-haslo', 'Hasło', 'current-password');
  const nowe = poleHasla('au-nowe', 'Nowe hasło', 'new-password');
  const powtorzenie = poleHasla('au-powtorzenie', 'Powtórz hasło', 'new-password');
  // Login, adres i droga potwierdzenia stoją w jednym formularzu i chowają się trybem.
  const login = poleTekstu('au-login', 'Login', 'username');
  const adres = poleTekstu('au-adres', 'Adres e-mail', 'email', 'email');
  const droga = poleTekstu('au-droga', 'Droga potwierdzenia z listu', 'one-time-code');

  const objasnienie = document.createElement('p');
  objasnienie.className = 'dn-pole-opis au-objasnienie';

  const przycisk = document.createElement('button');
  przycisk.type = 'submit';
  przycisk.className = 'dn-btn dn-btn--atrament au-wyslij';

  /** Pole „nie wyloguj mnie” rozstrzyga, czy sesja przeżyje zamknięcie aplikacji. */
  const niewylogowuj = utworzNiewylogowuj(sesjaTrwala());

  /** Segmenty metody wejścia obsadza wynik rozpoznania; przełączenie niczego nie wysyła ani nie blokuje. */
  const segmenty = utworzSegmentyMetody((wybrana) => wybierzMetode(wybrana));

  /** Odnośnik resetu przełącza ekran na odzyskanie konta i wraca; stoi wyłącznie przy wejściu hasłem. */
  const reset = utworzOdnosnikResetu(() =>
    pokazFormularz(tryb === 'wejscie' ? 'odzyskanie' : 'wejscie'),
  );

  /** Odnośnik wprowadza krok potwierdzenia adresu drogą z listu i pozwala z niego wyjść. */
  const potwierdzenieZListu = utworzOdnosnikPotwierdzenia(() =>
    pokazFormularz(tryb === 'potwierdzenie' ? 'wejscie' : 'potwierdzenie'),
  );

  const formularz = document.createElement('form');
  formularz.className = 'au-formularz';
  formularz.append(
    segmenty.element,
    login.pole,
    adres.pole,
    droga.pole,
    haslo.pole,
    nowe.pole,
    powtorzenie.pole,
    objasnienie,
    niewylogowuj.pole,
    przycisk,
    reset,
    potwierdzenieZListu,
  );
  formularz.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    void wyslij();
  });

  const ponow = document.createElement('button');
  ponow.type = 'button';
  ponow.className = 'dn-btn dn-btn--zarys au-ponow';
  ponow.textContent = 'Spróbuj ponownie';
  ponow.hidden = true;
  ponow.addEventListener('click', () => void rozpoznaj());

  const element = zlozPrzeslone(stany.element, ponow);

  // ── przepływ ───────────────────────────────────────────────────────────────

  /** Rozpoznaje stan bramki i pokazuje jego wynik; kroki rozpoznania prowadzi rdzeń. */
  async function rozpoznaj(): Promise<void> {
    ponow.hidden = true;
    stany.ladowanie('Sprawdzam stan bramki…');
    const droga = await rozpoznajWejscie({
      zrodlo,
      zapis: { odczytaj: odczytajSesje, skasuj: skasujSesje },
      urzadzenie: odczytajTozsamoscUrzadzenia,
      opiszOdmowe: (obszar, blad, bez) =>
        opisOdmowyBramki('Przedłużenie poprzedniej sesji', obszar, blad, bez),
    });

    if (droga.rodzaj === 'bez-przeslony') {
      // Przesłona schodzi bez komunikatu — nastawa wyłączająca logowanie jest świadoma.
      zdejmij();
      return;
    }
    if (droga.rodzaj === 'sesja') {
      wpusc(droga.sesja);
      return;
    }
    if (droga.rodzaj === 'niepewny') {
      stany.blad(
        opisOdmowyBramki('Rozpoznanie stanu bramki', 'rozpoznanie', droga.blad, droga.bezUchwytu),
        droga.blad,
      );
      ponow.hidden = false;
      return;
    }
    metody = droga.metody;
    pokazFormularz(droga.tryb, metody[0] ?? 'haslo');
    if (droga.notatka !== undefined) stany.potwierdzenie(droga.notatka, false);
  }

  function pokazFormularz(nowy: Tryb, nowaMetoda: MetodaWejscia = metoda): void {
    tryb = nowy;
    // Metoda inna niż hasło należy wyłącznie do wejścia; zmiana hasła używa hasła dotychczasowego.
    metoda = tryb === 'wejscie' ? nowaMetoda : 'haslo';
    segmenty.element.hidden = tryb !== 'wejscie';
    if (tryb === 'wejscie') segmenty.pokaz(metody, metoda);
    // Login idzie z rejestracją i wejściem hasłem; rdzeń rozpoznaje konto po nim.
    login.pole.hidden = !(tryb === 'zalozenie' || (tryb === 'wejscie' && metoda === 'haslo'));
    adres.pole.hidden = !(tryb === 'zalozenie' || tryb === 'odzyskanie');
    droga.pole.hidden = !(tryb === 'potwierdzenie' || tryb === 'odzyskanie-haslo');
    // Pole hasła dotychczasowego znika tam, gdzie z definicji hasła jeszcze nie ma.
    haslo.pole.hidden =
      tryb === 'potwierdzenie' || tryb === 'odzyskanie' || tryb === 'odzyskanie-haslo';
    // Pole „nowe hasło" należy do zmiany hasła i do ostatniego kroku odzyskania.
    nowe.pole.hidden = !(tryb === 'reset' || tryb === 'odzyskanie-haslo');
    powtorzenie.pole.hidden = tryb === 'wejscie' || tryb === 'potwierdzenie' || tryb === 'odzyskanie';
    haslo.ustawEtykiete(etykietaPierwszego(tryb, metoda));
    haslo.kontrolka.autocomplete = tryb === 'zalozenie' ? 'new-password' : 'current-password';
    // Puste objaśnienie znika, by nie zostawiać pustego odstępu pod polami formularza.
    objasnienie.textContent = OBJASNIENIA[tryb];
    objasnienie.hidden = OBJASNIENIA[tryb] === '';
    przycisk.textContent = NAPISY_PRZYCISKU[tryb];
    // Nastawa „nie wyloguj mnie” stoi tam, gdzie czynność kończy się wejściem do aplikacji.
    niewylogowuj.pole.hidden = tryb !== 'wejscie' && tryb !== 'potwierdzenie';
    reset.hidden = tryb === 'zalozenie';
    // Droga z listu stoi przy wejściu i przy potwierdzeniu — drogach powrotnych do konta.
    potwierdzenieZListu.hidden = tryb !== 'wejscie' && tryb !== 'potwierdzenie';
    if (tryb === 'potwierdzenie') potwierdzenieZListu.textContent = 'Wróć do logowania';
    else potwierdzenieZListu.textContent = 'Mam drogę potwierdzenia z listu';
    stany.tresc().replaceChildren(formularz);
    // Ognisko trafia na pierwsze pole widoczne w trybie, nigdy na pole ukryte.
    pierwszeWidoczne().focus();
  }

  /** Pierwsze pole widoczne w bieżącym trybie — cel ogniska po przełączeniu. */
  function pierwszeWidoczne(): HTMLInputElement {
    for (const kandydat of [login, adres, droga, haslo, nowe]) {
      if (!kandydat.pole.hidden) return kandydat.kontrolka;
    }
    return haslo.kontrolka;
  }

  /** Przełączenie metody czyści pole sekretu, by PIN nie trafił do rdzenia jako hasło. */
  function wybierzMetode(wybrana: MetodaWejscia): void {
    if (wybrana === metoda) return;
    haslo.kontrolka.value = '';
    pokazFormularz(tryb, wybrana);
  }

  /** Wysłanie formularza — rdzeń rozstrzyga, klient pokazuje jego odpowiedź. */
  async function wyslij(): Promise<void> {
    if (tryb === 'reset') {
      await zmienHaslo();
      return;
    }
    if (tryb === 'potwierdzenie') {
      await potwierdzAdres();
      return;
    }
    if (tryb === 'odzyskanie') {
      await poprosOOdzyskanie();
      return;
    }
    if (tryb === 'odzyskanie-haslo') {
      await ustawHasloPoOdzyskaniu();
      return;
    }
    // Ekran nie sprawdza długości PIN-u — minimum znaków dotyczy wyłącznie hasła.
    if (tryb === 'zalozenie' && haslo.kontrolka.value.length < MIN_ZNAKOW) {
      stany.potwierdzenie(
        `Hasło jest za krótkie — potrzeba co najmniej ${MIN_ZNAKOW} znaków. Nic nie zostało wysłane.`,
        false,
      );
      return;
    }
    if (tryb === 'zalozenie' && haslo.kontrolka.value !== powtorzenie.kontrolka.value) {
      // Zgodność hasła z powtórzeniem jest jedynym sprawdzianem formularza przy zakładaniu.
      stany.potwierdzenie(
        'Hasła różnią się od siebie — wpisz to samo hasło w obu polach. Nic nie zostało wysłane.',
        false,
      );
      return;
    }
    przycisk.setAttribute('aria-busy', 'true');
    try {
      if (tryb === 'wejscie') await wejdz();
      else await zaloz();
    } finally {
      przycisk.removeAttribute('aria-busy');
    }
  }

  async function wejdz(): Promise<void> {
    stany.potwierdzenie('Pytam rdzeń o wejście (auth.login)…', true);
    const odpowiedz = await zrodlo.wejdz({
      metoda,
      login: login.kontrolka.value.trim(),
      sekret: haslo.kontrolka.value,
      niewylogowuj: niewylogowuj.kontrolka.checked,
      urzadzenie: odczytajTozsamoscUrzadzenia(),
    });
    if (odpowiedz.udana && odpowiedz.wynik !== undefined) {
      wpusc(odpowiedz.wynik.session);
      return;
    }
    // Obszar rady zależy od metody — jedna rada na oba znaczenia myliłaby przy PIN-ie.
    const zdanie = opisOdmowyBramki(
      metoda === 'pin' ? 'Wejście PIN-em' : 'Wejście przez bramkę',
      metoda === 'pin' ? 'wejscie-pin' : 'wejscie',
      odpowiedz.blad,
      bezUchwytu(odpowiedz),
    );
    // Gdy rdzeń mówi „nieustawiona”, ekran przechodzi na założenie konta zamiast wejścia.
    if (metoda === 'haslo' && !bezUchwytu(odpowiedz) && odpowiedz.blad?.code === 'not_found') {
      pokazFormularz('zalozenie');
    }
    stany.potwierdzenie(zdanie, false);
  }

  /** Zmiana hasła nie wpuszcza — rdzeń oddaje potwierdzenie, a wejście wymaga nowego hasła. */
  async function zmienHaslo(): Promise<void> {
    if (nowe.kontrolka.value.length < MIN_ZNAKOW) {
      stany.potwierdzenie(
        `Nowe hasło jest za krótkie — potrzeba co najmniej ${MIN_ZNAKOW} znaków. Nic nie zostało wysłane.`,
        false,
      );
      return;
    }
    if (nowe.kontrolka.value !== powtorzenie.kontrolka.value) {
      stany.potwierdzenie(
        'Nowe hasła różnią się od siebie — wpisz to samo w obu polach. Nic nie zostało wysłane.',
        false,
      );
      return;
    }
    stany.potwierdzenie('Zmieniam hasło…', true);
    const odpowiedz = await zrodlo.zmienHaslo(haslo.kontrolka.value, nowe.kontrolka.value);
    if (!odpowiedz.udana || odpowiedz.wynik === undefined) {
      stany.potwierdzenie(
        opisOdmowyBramki('Zmiana hasła', 'zmiana', odpowiedz.blad, bezUchwytu(odpowiedz)),
        false,
      );
      return;
    }
    // Zmiana unieważnia zalogowania — także to zapisane na tej maszynie.
    skasujSesje();
    haslo.kontrolka.value = '';
    nowe.kontrolka.value = '';
    powtorzenie.kontrolka.value = '';
    pokazFormularz('wejscie');
    stany.potwierdzenie(
      odpowiedz.wynik.changed
        ? 'Hasło zmienione. Zaloguj się nowym hasłem.'
        : 'Rdzeń przyjął wywołanie, ale hasła nie zmienił.',
      odpowiedz.wynik.changed,
    );
  }

  async function zaloz(): Promise<void> {
    stany.potwierdzenie('Zakładam konto…', true);
    const odpowiedz = await zrodlo.zaloz({
      login: login.kontrolka.value.trim(),
      adres: adres.kontrolka.value.trim(),
      haslo: haslo.kontrolka.value,
    });
    if (odpowiedz.udana && odpowiedz.wynik !== undefined) {
      // Hasło i powtórzenie znikają z pól po założeniu konta, przed potwierdzeniem adresu.
      haslo.kontrolka.value = '';
      powtorzenie.kontrolka.value = '';
      pokazFormularz('potwierdzenie');
      stany.potwierdzenie(
        `Konto założone. List z drogą potwierdzenia poszedł na ${adres.kontrolka.value.trim()}.`,
        true,
      );
      return;
    }
    const zdanie = opisOdmowyBramki(
      'Założenie konta',
      'zalozenie',
      odpowiedz.blad,
      bezUchwytu(odpowiedz),
    );
    // Gdy konto już istnieje, jedyną drogą dalej jest wejście hasłem istniejącym.
    if (!bezUchwytu(odpowiedz) && odpowiedz.blad?.code === 'conflict') {
      pokazFormularz('wejscie');
    }
    stany.potwierdzenie(zdanie, false);
  }

  /** Krok drugi rejestracji: rdzeń wydaje token dostępu dopiero po potwierdzeniu adresu. */
  async function potwierdzAdres(): Promise<void> {
    stany.potwierdzenie('Potwierdzam adres…', true);
    const odpowiedz = await zrodlo.potwierdz(
      droga.kontrolka.value.trim(),
      niewylogowuj.kontrolka.checked,
    );
    if (odpowiedz.udana && odpowiedz.wynik !== undefined) {
      droga.kontrolka.value = '';
      wpusc(odpowiedz.wynik.session);
      return;
    }
    stany.potwierdzenie(
      opisOdmowyBramki('Potwierdzenie adresu', 'potwierdzenie', odpowiedz.blad, bezUchwytu(odpowiedz)),
      false,
    );
  }

  /** Odpowiedź rdzenia jest jednakowa dla adresu istniejącego i obcego konta. */
  async function poprosOOdzyskanie(): Promise<void> {
    stany.potwierdzenie('Wysyłam drogę odzyskania…', true);
    const odpowiedz = await zrodlo.odzyskaj(adres.kontrolka.value.trim());
    if (odpowiedz.udana && odpowiedz.wynik !== undefined) {
      pokazFormularz('odzyskanie-haslo');
      stany.potwierdzenie(
        'Jeżeli ten adres jest adresem uwierzytelniającym konta, poszła na niego droga ' +
          'potwierdzenia. Przepisz ją poniżej i podaj nowe hasło.',
        true,
      );
      return;
    }
    stany.potwierdzenie(
      opisOdmowyBramki('Odzyskanie konta', 'odzyskanie', odpowiedz.blad, bezUchwytu(odpowiedz)),
      false,
    );
  }

  /** Odzyskanie konta, krok drugi — ustawienie nowego hasła drogą z listu. */
  async function ustawHasloPoOdzyskaniu(): Promise<void> {
    if (nowe.kontrolka.value.length < MIN_ZNAKOW) {
      stany.potwierdzenie(
        `Hasło jest za krótkie — potrzeba co najmniej ${MIN_ZNAKOW} znaków. Nic nie zostało wysłane.`,
        false,
      );
      return;
    }
    if (nowe.kontrolka.value !== powtorzenie.kontrolka.value) {
      stany.potwierdzenie(
        'Nowe hasła różnią się od siebie — wpisz to samo w obu polach. Nic nie zostało wysłane.',
        false,
      );
      return;
    }
    stany.potwierdzenie('Ustawiam nowe hasło…', true);
    const odpowiedz = await zrodlo.ustawNoweHaslo(
      droga.kontrolka.value.trim(),
      nowe.kontrolka.value,
    );
    if (!odpowiedz.udana || odpowiedz.wynik === undefined) {
      stany.potwierdzenie(
        opisOdmowyBramki('Ustawienie nowego hasła', 'odzyskanie', odpowiedz.blad, bezUchwytu(odpowiedz)),
        false,
      );
      return;
    }
    // Odzyskanie unieważnia wszystkie zalogowania — także zapis na tej maszynie.
    skasujSesje();
    droga.kontrolka.value = '';
    nowe.kontrolka.value = '';
    powtorzenie.kontrolka.value = '';
    pokazFormularz('wejscie');
    stany.potwierdzenie(
      odpowiedz.wynik.changed
        ? `Hasło ustawione. Zamknięto zalogowania na ${odpowiedz.wynik.revokedDevices} urządzeniach — zaloguj się nowym hasłem.`
        : 'Rdzeń przyjął wywołanie, ale hasła nie zmienił.',
      odpowiedz.wynik.changed,
    );
  }

  /** Wejście: zapis sesji, oddanie jej aplikacji, zdjęcie przesłony — bez klikania. */
  function wpusc(sesja: AuthSession): void {
    if (wpuszczony) return;
    wpuszczony = true;
    zapiszSesje(sesja, niewylogowuj.kontrolka.checked);
    stany.potwierdzenie(`Wejście otwarte — ${opisWaznosci(sesja.expiresAt)}.`, true);
    opis.naWejscie(sesja);
    element.remove();
  }

  /** Zdjęcie przesłony bez wejścia — aplikacja pod spodem jest już złożona i połączona. */
  function zdejmij(): void {
    if (wpuszczony) return;
    wpuszczony = true;
    element.remove();
  }

  return {
    element,

    uruchom() {
      if (wpuszczony) return;
      if (!element.isConnected) document.body.append(element);
      void rozpoznaj();
    },

    rozlacz() {
      element.remove();
    },
  };
}
