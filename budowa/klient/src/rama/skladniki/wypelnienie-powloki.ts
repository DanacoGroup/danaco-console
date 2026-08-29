/**
 * Wypełnienie powłoki wartościami rdzenia.
 *
 * Powłoka pochodzi z biblioteki prototypu (`zasoby/powloka.js`) i przychodzi
 * z treścią przykładową: nazwą okna rejestracji w belce, adresem
 * `login@przykład.pl` w pasku stanu oraz miarami `CPU 18%` i `RAM 2,4 GB`.
 * W prototypie były po to, żeby dało się ocenić układ; w produkcie byłyby
 * atrapą podaną Operatorowi jako pomiar.
 *
 * Ten plik zastępuje je wartościami, które rdzeń naprawdę oddaje, a pozycje
 * bez źródła zdejmuje z paska. Zdjęcie jest świadome: miary obciążenia maszyny
 * nie stoją w kontrakcie, więc nie ma ich skąd wziąć, a zostawione kłamałyby.
 */

import { tekst } from '../narzedzia.ts';

export interface WartosciPowloki {
  /** Nazwa okna widoczna w belce tytułowej. */
  tytul: string;
  /** Środowisko, do którego przebieg wszedł. */
  srodowisko: string;
  /** Liczba kart sesji odtworzonych przez rdzeń. */
  liczbaSesji: number;
  /** Motyw rozstrzygnięty przez wybór Operatora albo ustawienie systemu. */
  motywCiemny: boolean;
  /** Konto Operatora; puste znaczy konto nierozpoznane. */
  operator?: string;
}

/**
 * Pozycje paska stanu bez pokrycia w kontrakcie. Zdejmowane po treści, bo
 * biblioteka nie znaczy ich atrybutem — a treść przykładowa jest tu jedynym
 * pewnym rozpoznaniem.
 */
const BEZ_ZRODLA = ['CPU', 'RAM'];

function ustawTresc(korzen: ParentNode, wybor: string, wartosc: string): void {
  const wezel = korzen.querySelector(wybor);
  if (wezel !== null) wezel.textContent = wartosc;
}

export function wypelnijPowloke(korzen: ParentNode, w: WartosciPowloki): void {
  ustawTresc(korzen, '.dn-belka-tytul', w.tytul);

  /* Prototyp Studia nakłada na prawą kolumnę `st-obudowa` — to ona daje jej
     układ okna roboczego. Biblioteka powłoki jest wspólna dla wszystkich okien,
     więc klasy właściwej oknu nie niesie. */
  korzen.querySelector('.dn-rama-prawa')?.classList.add('st-obudowa');

  const pozycje = [...korzen.querySelectorAll('.dn-stan-poz')];
  for (const pozycja of pozycje) {
    const tresc = pozycja.textContent ?? '';

    /* Miary maszyny zdejmowane w całości: pozycja pusta czytałaby się jak
       pomiar zerowy, a zdanie „brak pomiaru” zajmowałoby miejsce niczym. */
    if (BEZ_ZRODLA.some((nazwa) => tresc.includes(nazwa))) {
      pozycja.remove();
      continue;
    }

    if (tresc.includes(tekst('stan.operator'))) {
      ustawTresc(pozycja, 'span:last-child', '');
      pozycja.textContent = `${tekst('stan.operator')} · ${w.operator ?? tekst('stan.brakTozsamosci')}`;
      continue;
    }
    if (tresc.includes(tekst('stan.sesje'))) {
      pozycja.textContent = `${tekst('stan.sesje')}: ${w.liczbaSesji}`;
      continue;
    }
    if (tresc.includes(tekst('stan.widok'))) {
      const motyw = w.motywCiemny ? tekst('stan.motywCiemny') : tekst('stan.motywJasny');
      pozycja.textContent = `${tekst('stan.widok')} ${motyw}`;
      continue;
    }
    /* Pozycja „— · —” to miejsce na środowisko: prototyp zostawił w niej
       kreski, bo w chwili rysowania środowiska jeszcze nie było. */
    if (tresc.replaceAll(/[\s·—-]/gu, '') === '') {
      pozycja.textContent = `${tekst('stan.srodowisko')} · ${w.srodowisko}`;
    }
  }
}
