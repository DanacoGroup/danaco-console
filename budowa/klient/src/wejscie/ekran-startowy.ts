/**
 * Wpięcie ekranu startowego — znaku marki kreślonego światłem, który gra po
 * uruchomieniu programu, zanim stanie okno drogi wejścia.
 *
 * Animacji nie odtwarza ten plik: niesie ją składnik biblioteki
 * `design/zasoby/ekran-startowy.js`, wciągnięty przez dokument. Tutaj stoi
 * wyłącznie sterowanie z umowy tego składnika — odtworzenie, zgłoszenie
 * gotowości programu i sprzątnięcie warstwy po jej wygaszeniu.
 *
 * `data-czekaj` na polu znaczy, że składnik zatrzyma się przed wygaszeniem
 * i poczeka na `gotowe()`. Dzięki temu animacja nie gaśnie przed oknem, gdy
 * połączenie z serwerem idzie szybciej niż trzy sekundy jej biegu, i nie
 * zostaje na ekranie, gdy idzie wolniej.
 */

/** Sterowanie składnika biblioteki, wystawione na węźle pola. */
interface SterowanieEkranu {
  odtworz(nastawy?: { poKoncu?: () => void; czekajNaGotowosc?: boolean }): void;
  gotowe(): void;
  pomin(): void;
  zdejmij(): void;
}

type PoleEkranu = HTMLElement & { ekranStartowy?: SterowanieEkranu };

/** Najdłuższe czekanie warstwy na gotowość przebiegu. Animacja trwa 3 s; ta granica daje jej zapas i zamyka drogę do zawieszenia okna. */
const GRANICA_CZEKANIA_MS = 6000;

export interface ZamontowanyEkranStartowy {
  /** Program wstał — wolno domknąć animację i odsłonić okno. */
  gotowe(): void;
}

/**
 * Uruchamia ekran startowy w dokumencie. Brak pola albo brak składnika
 * biblioteki nie jest usterką zatrzymującą program: droga wejścia działa bez
 * animacji, więc warstwa zostaje wtedy zdjęta od razu, a okno staje bez zwłoki.
 */
export function zalozEkranStartowy(dokument: Document): ZamontowanyEkranStartowy {
  const pole = dokument.querySelector<PoleEkranu>('[data-ekran-startowy]');
  const warstwa = dokument.querySelector<HTMLElement>('[data-uruchomienie]');
  /* Znacznik niesie każda scena, nie jedna: reguła biblioteki wiąże ukrycie
     okna ze sceną, a przebieg przechodzi między nimi w trakcie animacji. */
  const sceny = [...dokument.querySelectorAll<HTMLElement>('.we-scena')];

  function odsloniecie(): void {
    for (const scena of sceny) scena.removeAttribute('data-uruchamianie');
    warstwa?.remove();
  }

  for (const scena of sceny) scena.setAttribute('data-uruchamianie', '');

  if (pole === null) {
    odsloniecie();
    return { gotowe: () => undefined };
  }

  /* Składnik biblioteki zakłada się na polu dopiero przy `DOMContentLoaded`,
     a ten moduł biegnie wcześniej — jest odroczony, więc wykonuje się przed tym
     zdarzeniem. Odczyt sterowania od razu zastawał puste pole i zdejmował
     warstwę, zanim animacja miała szansę ruszyć. */
  let sterowanie: SterowanieEkranu | undefined;
  let gotowoscZgloszona = false;

  function uruchom(): void {
    sterowanie = pole?.ekranStartowy;
    if (sterowanie === undefined) {
      odsloniecie();
      return;
    }
    sterowanie.odtworz({ poKoncu: odsloniecie, czekajNaGotowosc: true });
    if (gotowoscZgloszona) sterowanie.gotowe();
    /* Bezpiecznik: `data-czekaj` trzyma animację do zgłoszenia gotowości, więc
       przebieg, który nie odpowie ani razu — zerwane gniazdo, odmowa pochodzenia
       — zostawiłby Operatora przed samym znakiem marki bez żadnego wyjścia.
       Po tym czasie warstwa ustępuje sama, a droga wejścia pokazuje swój
       własny stan błędu, który Operator potrafi przeczytać. */
    globalThis.setTimeout(() => {
      if (warstwa?.isConnected === true) sterowanie?.pomin();
    }, GRANICA_CZEKANIA_MS);
  }

  /* Czekamy do `complete`, nie do samego `loading`: moduł odroczony wykonuje się
     przy `interactive`, a składnik biblioteki zakłada się dopiero w obsłudze
     `DOMContentLoaded`, czyli po nim. Warunek na `loading` przepuszczał ten
     moduł przodem i warstwa znikała, zanim animacja ruszyła. */
  if (dokument.readyState === 'complete') {
    uruchom();
  } else {
    dokument.addEventListener('DOMContentLoaded', uruchom, { once: true });
  }

  return {
    gotowe() {
      gotowoscZgloszona = true;
      sterowanie?.gotowe();
    },
  };
}
