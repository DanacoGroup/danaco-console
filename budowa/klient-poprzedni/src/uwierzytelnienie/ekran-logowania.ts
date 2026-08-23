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

/**
 * Ekran logowania — jedyna bramka produktu; po wejściu nic więcej nie pyta.
 *
 * Ekran jest przesłoną nad aplikacją, nie czwartą trasą: kładzie się nad
 * gospodarzem dokumentu, a aplikacja pod nim składa się i łączy z rdzeniem w tym
 * samym czasie, więc po wejściu przesłona znika i widoczne jest gotowe Centrum
 * dowodzenia, zamiast drugiego ładowania. Trasy (`aplikacja/trasy.ts`) zostają
 * nietknięte — bramka nie jest miejscem pracy, więc nie jest trasą.
 *
 * Ekran nie zakłada budzika, nie liczy czasu sesji i nie przerywa pracy pytaniem
 * o tożsamość. Rdzeń nie odcina komend po wygaśnięciu sesji — bramka jest
 * progiem wejścia, nie strażnikiem każdego żądania — więc wygaśnięcie ujawnia
 * się wyłącznie tutaj, przy następnym uruchomieniu, zdaniem nad formularzem.
 *
 * Który z dwóch formularzy pokazać — wejście hasłem czy pierwsze ustawienie
 * hasła — rozstrzyga rdzeń (`rozpoznanie-bramki.ts`), nie domysł klienta. Pole
 * nie jest blokowane, przycisk nie jest wyszarzany; jedyny sprawdzian po stronie
 * formularza to zgodność hasła z powtórzeniem przy zakładaniu, bo kontrakt
 * opisuje powtórzenie jako sprawę formularza klienta.
 *
 * Wymóg logowania jest nastawą (`gateway.requireLogin`) — gdy jest wyłączony,
 * przesłona nie staje. Kroki rozpoznania: `rozpoznanie-bramki.ts`.
 */
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
  // Login i adres należą do rejestracji i do wejścia hasłem; droga potwierdzenia
  // — do dwóch kroków, które przychodzą listem. Wszystkie trzy stoją w jednym
  // formularzu i chowają się trybem, zamiast budować trzy osobne postacie ekranu.
  const login = poleTekstu('au-login', 'Login', 'username');
  const adres = poleTekstu('au-adres', 'Adres e-mail', 'email', 'email');
  const droga = poleTekstu('au-droga', 'Droga potwierdzenia z listu', 'one-time-code');

  const objasnienie = document.createElement('p');
  objasnienie.className = 'dn-pole-opis au-objasnienie';

  const przycisk = document.createElement('button');
  przycisk.type = 'submit';
  przycisk.className = 'dn-btn dn-btn--atrament au-wyslij';

  /**
   * „Nie wyloguj mnie" — rozstrzyga, czy sesja przeżyje zamknięcie aplikacji.
   *
   * Stan wyjściowy bierze się z tego, gdzie leży sesja zapisana poprzednio
   * (`sesjaTrwala`), a nie ze stałej: raz odznaczone pole ma zostać odznaczone.
   * Skutek opisuje `sesja-bramki.ts` — zaznaczone kładzie sesję w pamięci
   * trwałej, odznaczone w pamięci okna — a ta sama wartość idzie do rdzenia
   * i rozstrzyga o trwaniu sesji.
   */
  const niewylogowuj = utworzNiewylogowuj(sesjaTrwala());

  /**
   * Segmenty metody wejścia — obsadzane wynikiem rozpoznania, nie na zapas.
   * Przełączenie niczego nie wysyła i niczego nie blokuje (patrz `wybierzMetode`).
   */
  const segmenty = utworzSegmentyMetody((wybrana) => wybierzMetode(wybrana));

  /**
   * „Reset hasła" — przestawia ekran na zmianę hasła i z powrotem.
   *
   * Stoi wyłącznie przy wejściu hasłem. Przy pierwszym uruchomieniu nie ma
   * czego resetować (hasła jeszcze nie ma), a w samym trybie zmiany odnośnik
   * prowadziłby tam, gdzie ekran już stoi.
   */
  // Odnośnik prowadzi do ODZYSKANIA konta, nie do zmiany hasła ze znanym hasłem
  // bieżącym: naciska go ten, kto hasła nie pamięta, a wtedy zmiana ze znanym
  // hasłem jest drogą donikąd. Zmiana hasła świadoma jest czynnością Ustawień.
  const reset = utworzOdnosnikResetu(() =>
    pokazFormularz(tryb === 'wejscie' ? 'odzyskanie' : 'wejscie'),
  );

  /**
   * „Mam drogę potwierdzenia z listu" — wejście w krok drugi rejestracji
   * i wyjście z niego.
   *
   * Bez tego odnośnika krok potwierdzenia prowadził wyłącznie z udanej
   * rejestracji w tym samym oknie. Odświeżenie strony między listem
   * a przepisaniem drogi zostawiało Operatora przed formularzem wejścia bez
   * żadnej drogi dalej: rejestracja odmawia wtedy `conflict`, a logowanie —
   * oczekiwaniem na potwierdzenie adresu. To samo dotyczy listu odczytanego
   * na innej maszynie.
   */
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

  /**
   * Rozpoznanie: co pokazać, zanim ekran cokolwiek wyświetli.
   *
   * Kolejność kroków i powód każdego z nich stoją w `rozpoznanie-bramki.ts` —
   * tutaj zostaje wyłącznie pokazanie wyniku.
   */
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
      // Przesłona schodzi bez komunikatu: zdanie „logowanie wyłączone" byłoby
      // meldunkiem o nastawie ustawionej świadomie, a aplikacja pod spodem jest
      // już złożona i gotowa.
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
    // Metoda inna niż hasło należy wyłącznie do wejścia: pierwsze hasło zakłada
    // się hasłem, a zmienia hasłem dotychczasowym.
    metoda = tryb === 'wejscie' ? nowaMetoda : 'haslo';
    segmenty.element.hidden = tryb !== 'wejscie';
    if (tryb === 'wejscie') segmenty.pokaz(metody, metoda);
    // Login idzie z rejestracją i z wejściem hasłem — rdzeń rozpoznaje konto po
    // nim. Metody właściwe urządzeniu loginu nie potrzebują, bo wskazuje je
    // materiał leżący na maszynie.
    login.pole.hidden = !(tryb === 'zalozenie' || (tryb === 'wejscie' && metoda === 'haslo'));
    adres.pole.hidden = !(tryb === 'zalozenie' || tryb === 'odzyskanie');
    droga.pole.hidden = !(tryb === 'potwierdzenie' || tryb === 'odzyskanie-haslo');
    // Pierwszego pola nie ma tam, gdzie nie pyta się o hasło dotychczasowe:
    // przy potwierdzeniu adresu, przy prośbie o odzyskanie i przy ustawianiu
    // nowego hasła po odzyskaniu — tam hasła dotychczasowego z definicji nie ma,
    // bo właśnie dlatego konto się odzyskuje.
    haslo.pole.hidden =
      tryb === 'potwierdzenie' || tryb === 'odzyskanie' || tryb === 'odzyskanie-haslo';
    // Pole „nowe hasło" należy do zmiany hasła i do ostatniego kroku odzyskania.
    nowe.pole.hidden = !(tryb === 'reset' || tryb === 'odzyskanie-haslo');
    powtorzenie.pole.hidden = tryb === 'wejscie' || tryb === 'potwierdzenie' || tryb === 'odzyskanie';
    haslo.ustawEtykiete(etykietaPierwszego(tryb, metoda));
    haslo.kontrolka.autocomplete = tryb === 'zalozenie' ? 'new-password' : 'current-password';
    // Puste objaśnienie znika, zamiast zostawiać pusty odstęp pod polami:
    // przy wejściu hasłem ekran nie ma nic do dopowiedzenia.
    objasnienie.textContent = OBJASNIENIA[tryb];
    objasnienie.hidden = OBJASNIENIA[tryb] === '';
    przycisk.textContent = NAPISY_PRZYCISKU[tryb];
    // Nastawa „nie wyloguj mnie" stoi tam, gdzie czynność kończy się wejściem:
    // przy logowaniu i przy potwierdzeniu adresu, które wydaje pierwszą sesję.
    // Rejestracja sesji nie zakłada, więc nie ma tam czego przedłużać.
    niewylogowuj.pole.hidden = tryb !== 'wejscie' && tryb !== 'potwierdzenie';
    reset.hidden = tryb === 'zalozenie';
    // Droga z listu stoi przy wejściu — tam wraca ten, kto konto założył i list
    // odczytał później — oraz w samym kroku potwierdzenia, gdzie jest drogą
    // powrotną. Przy zakładaniu konta listu jeszcze nie ma, a przy odzyskaniu
    // drogę niesie już własne pole tamtych kroków.
    potwierdzenieZListu.hidden = tryb !== 'wejscie' && tryb !== 'potwierdzenie';
    if (tryb === 'potwierdzenie') potwierdzenieZListu.textContent = 'Wróć do logowania';
    else potwierdzenieZListu.textContent = 'Mam drogę potwierdzenia z listu';
    stany.tresc().replaceChildren(formularz);
    // Ognisko idzie na pierwsze pole widoczne w tym trybie, nie zawsze na hasło:
    // w kroku potwierdzenia hasła nie ma, a ognisko na polu ukrytym zostawiłoby
    // Operatora przed formularzem, w którym nic się nie dzieje po naciśnięciu
    // klawisza.
    pierwszeWidoczne().focus();
  }

  /** Pierwsze pole widoczne w bieżącym trybie — cel ogniska po przełączeniu. */
  function pierwszeWidoczne(): HTMLInputElement {
    for (const kandydat of [login, adres, droga, haslo, nowe]) {
      if (!kandydat.pole.hidden) return kandydat.kontrolka;
    }
    return haslo.kontrolka;
  }

  /**
   * Przełączenie metody wejścia.
   *
   * Sekret znika przy przełączeniu: PIN wpisany w polu podpisanym „Hasło"
   * poleciałby do rdzenia jako hasło, podniósł dławik prób i wrócił odmową
   * o haśle, którego nikt nie próbował podać.
   */
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
    // Długości PIN-u ekran nie pilnuje. Minimum znaków jest regułą formularza
    // dla hasła zakładanego tutaj; PIN zakłada się w Ustawieniach i to tamten
    // formularz obiecuje jego kształt. Sprawdzian przepisany na wejście PIN-em
    // odmawiałby PIN-owi, który rdzeń przyjmie.
    if (tryb === 'zalozenie' && haslo.kontrolka.value.length < MIN_ZNAKOW) {
      stany.potwierdzenie(
        `Hasło jest za krótkie — potrzeba co najmniej ${MIN_ZNAKOW} znaków. Nic nie zostało wysłane.`,
        false,
      );
      return;
    }
    if (tryb === 'zalozenie' && haslo.kontrolka.value !== powtorzenie.kontrolka.value) {
      // Jedyny sprawdzian formularza — nakazany przez kontrakt, bo rdzeń
      // dostaje hasło raz i przepisania sprawdzić nie może.
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
    // Obszar rady zależy od metody: `not_found` przy haśle znaczy „bramki nie
    // ustawiono", a przy PIN-ie „PIN-u na tej maszynie nie ma". Jedna rada na
    // dwa znaczenia odsyłałaby do zakładania hasła, które już stoi.
    const zdanie = opisOdmowyBramki(
      metoda === 'pin' ? 'Wejście PIN-em' : 'Wejście przez bramkę',
      metoda === 'pin' ? 'wejscie-pin' : 'wejscie',
      odpowiedz.blad,
      bezUchwytu(odpowiedz),
    );
    // Bramka zniknęła między rozpoznaniem a wejściem (świeża baza, inne
    // urządzenie): rdzeń mówi „nieustawiona", więc ekran przechodzi na
    // założenie. PIN-u to nie dotyczy — tam `not_found` mówi o samym PIN-ie.
    if (metoda === 'haslo' && !bezUchwytu(odpowiedz) && odpowiedz.blad?.code === 'not_found') {
      pokazFormularz('zalozenie');
    }
    stany.potwierdzenie(zdanie, false);
  }

  /**
   * Zmiana hasła z ekranu wejścia (`auth.password.reset`).
   *
   * Nie wpuszcza. Rdzeń oddaje tu potwierdzenie zmiany i liczbę unieważnionych
   * zalogowań, nie sesję — po zmianie trzeba zalogować się nowym hasłem. Ekran
   * wraca więc do formularza wejścia z wyczyszczonymi polami.
   */
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
      // Hasło i powtórzenie znikają z pól: konto jest założone, a droga
      // potwierdzenia przychodzi listem — trzymanie hasła w formularzu przez
      // czas czekania na list nie służy już niczemu.
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
    // Kotwica już stoi (założona z innego urządzenia albo wyścig dwóch okien):
    // jedyną drogą jest wejście hasłem istniejącym.
    if (!bezUchwytu(odpowiedz) && odpowiedz.blad?.code === 'conflict') {
      pokazFormularz('wejscie');
    }
    stany.potwierdzenie(zdanie, false);
  }

  /**
   * Krok drugi rejestracji — potwierdzenie adresu drogą z listu.
   *
   * To tutaj Operator wchodzi do platformy po raz pierwszy: rdzeń wydaje token
   * dostępu dopiero po potwierdzeniu.
   */
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

  /**
   * Odzyskanie konta, krok pierwszy — prośba o list.
   *
   * Odpowiedź rdzenia jest taka sama dla adresu właściciela i dla obcego, więc
   * ekran mówi dokładnie tyle, ile wie: że jeżeli adres pasuje, list poszedł.
   * Zdanie „wysłano" bez tego zastrzeżenia byłoby potwierdzeniem, że konto o tym
   * adresie istnieje.
   */
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

  /**
   * Zdjęcie przesłony bez wejścia — gdy strzec nie ma czego.
   *
   * `naWejscie` nie leci, bo nie ma sesji do oddania, a wywołanie go z sesją
   * zmyśloną byłoby atrapą wejścia. Aplikacja pod spodem jest już złożona
   * i połączona, więc zdjęcie przesłony jest całą czynnością.
   */
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
