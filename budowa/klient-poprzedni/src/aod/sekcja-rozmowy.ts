import { opisOdmowyAod } from './odmowy-aod';
import { utworzAkapit, utworzPodtytul } from './pola-wykazu';
import type { ZrodloAod } from './zrodlo-komend';

/**
 * Sekcja wysyłki wiadomości z nakładki do okna rozmowy (`aod.chat.send`).
 *
 * Pole okna zostaje puste z zamysłem: puste `windowId` kieruje wiadomość do
 * okna ogniskowanego, a ognisko zna rdzeń — mogło się zmienić między odczytem
 * stanu a wysyłką, więc sekcja nie podstawia tu okna odczytanego wcześniej.
 * Odpowiedź niesie `messageId` oraz `windowId` okna, które wiadomość przyjęło,
 * i meldunek pokazuje oba.
 *
 * Przycisk wysyłki pozostaje czynny także przy pustym polu treści — pustą treść
 * ocenia rdzeń i to on zwraca odmowę (`validation_failed`).
 */
export interface SekcjaRozmowy {
  element: HTMLElement;
  /** Wskazuje sesję, do której należy kierować wiadomość bez wskazanego okna. */
  ustawSesje(sessionId?: string): void;
}

export interface OpisSekcjiRozmowy {
  zrodlo: ZrodloAod;
  /** Krótkie potwierdzenie czynności na pasku okna. */
  zamelduj(zdanie: string, udane: boolean): void;
}

export function utworzSekcjeRozmowy(opis: OpisSekcjiRozmowy): SekcjaRozmowy {
  const element = document.createElement('section');
  element.className = 'ao-sekcja';

  const tresc = document.createElement('textarea');
  tresc.className = 'dn-pole-kontrolka ao-rozmowa__tresc';
  tresc.id = 'ao-rozmowa-tresc';
  tresc.rows = 3;
  tresc.placeholder = 'Treść wiadomości do okna rozmowy';

  const etykietaTresci = document.createElement('label');
  etykietaTresci.className = 'dn-pole-etykieta';
  etykietaTresci.htmlFor = tresc.id;
  etykietaTresci.textContent = 'Wiadomość';

  const okno = document.createElement('input');
  okno.type = 'text';
  okno.className = 'dn-pole-kontrolka ao-rozmowa__okno';
  okno.id = 'ao-rozmowa-okno';
  okno.placeholder = 'puste = okno ogniskowane';

  const etykietaOkna = document.createElement('label');
  etykietaOkna.className = 'dn-pole-etykieta';
  etykietaOkna.htmlFor = okno.id;
  etykietaOkna.textContent = 'Okno rozmowy';

  const wyslij = document.createElement('button');
  wyslij.type = 'submit';
  wyslij.className = 'dn-btn dn-btn--zarys';
  wyslij.textContent = 'Wyślij z nakładki';

  const formularz = document.createElement('form');
  formularz.className = 'ao-formularz ao-rozmowa';
  formularz.append(etykietaTresci, tresc, etykietaOkna, okno, wyslij);
  formularz.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    void wyslijWiadomosc();
  });

  element.append(
    utworzPodtytul('Rozmowa z nakładki (aod.chat.send)'),
    utworzAkapit(
      'ao-pusto',
      'Puste pole okna kieruje wiadomość do okna ogniskowanego — rdzeń odpowiada, ' +
        'które okno ją przyjęło.',
    ),
    formularz,
  );

  /** Sesja bieżąca — idzie z wiadomością, gdy okno nie zostało wskazane. */
  let idSesji: string | undefined;

  async function wyslijWiadomosc(): Promise<void> {
    const wskazaneOkno = okno.value.trim();

    const wynik = await opis.zrodlo.wyslijRozmowe({
      text: tresc.value,
      windowId: wskazaneOkno === '' ? undefined : wskazaneOkno,
      sessionId: idSesji,
    });

    if (!wynik.udany || wynik.wynik === undefined) {
      opis.zamelduj(opisOdmowyAod('Wysyłka wiadomości z nakładki', 'rozmowa', wynik.blad), false);
      return;
    }

    tresc.value = '';
    opis.zamelduj(
      `Wiadomość ${wynik.wynik.messageId} przyjęta przez okno ${wynik.wynik.windowId}.`,
      true,
    );
  }

  return {
    element,

    ustawSesje(sessionId) {
      idSesji = sessionId;
    },
  };
}
