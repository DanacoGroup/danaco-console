import type { AssistantVoiceCommandResponse } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';
import { zglosBrak } from './braki-kontraktu';

/**
 * Odpowiedź rdzenia na polecenie — transkrypcja, odnośnik syntezy i zamiana
 * fragmentu w zadanie. Panel pokazuje to, co wróciło z `assistant.voice.command`;
 * rdzenia nie wywołuje sam, a zamiary oddaje oknu.
 */
export interface PanelOdpowiedzi {
  element: HTMLElement;
  /** Nanosi odpowiedź rdzenia. */
  pokaz(odpowiedz: AssistantVoiceCommandResponse): void;
  /** Zdejmuje odpowiedź — polecenie wyszło, wyniku jeszcze nie ma. */
  wyczysc(): void;
}

/**
 * Zamiana zaznaczonego fragmentu odpowiedzi w zadanie Actions Monitora. Panel
 * podaje wyłącznie treść fragmentu, więc założenie zadania i wybór komendy
 * zostają po stronie okna.
 */
export type NaZadanie = (tresc: string) => void;

/**
 * Odsłuch odnośnika wykonywany przez okno komendą `speech.audio.fetch`. Panel
 * nie sięga po bajty sam, bo kanał do rdzenia trzyma okno, a panel odpowiada
 * wyłącznie za postać odpowiedzi.
 */
export type NaOdsluch = (odnosnik: string) => void;

export function utworzPanelOdpowiedzi(
  naZadanie: NaZadanie,
  naOdsluch: NaOdsluch,
): PanelOdpowiedzi {
  let ostatnia = '';

  const transkrypcja = document.createElement('p');
  transkrypcja.className = 'ma-odpowiedz__transkrypcja';

  const odnosnik = document.createElement('p');
  odnosnik.className = 'dn-pole-opis ma-odpowiedz__odnosnik';

  const odsluch = przycisk('Odsłuchaj odpowiedź', 'dn-btn dn-btn--sm dn-btn--zarys');
  let odnosnikSyntezy = '';
  odsluch.addEventListener('click', () => {
    if (odnosnikSyntezy === '') {
      zglosBrak(
        'Odsłuch syntezy',
        'Rdzeń nie oddał odnośnika odpowiedzi syntezowanej — polecenie poszło bez ' +
          'prośby o odczyt (przełącznik „Czytaj odpowiedź") albo synteza nie ruszyła.',
      );
      return;
    }
    naOdsluch(odnosnikSyntezy);
  });

  const zadanie = przycisk('Wstaw fragment jako zadanie', 'dn-btn dn-btn--sm dn-btn--zarys');
  zadanie.addEventListener('click', () => naZadanie(fragment(transkrypcja, ostatnia)));

  const przyciski = document.createElement('div');
  przyciski.className = 'ma-odpowiedz__przyciski';
  przyciski.append(odsluch, zadanie);

  const element = document.createElement('div');
  element.className = 'ma-odpowiedz';
  element.dataset['panel'] = 'odpowiedz';
  element.hidden = true;
  element.append(transkrypcja, odnosnik, przyciski);

  return {
    element,

    pokaz(odpowiedz) {
      ostatnia = odpowiedz.transcript;
      transkrypcja.textContent = `Rdzeń przyjął polecenie: „${odpowiedz.transcript}"`;
      odnosnikSyntezy = odpowiedz.speechRef ?? '';
      odnosnik.textContent =
        odnosnikSyntezy === ''
          ? 'Rdzeń nie oddał odnośnika odpowiedzi syntezowanej.'
          : `Odnośnik odpowiedzi syntezowanej: ${odnosnikSyntezy}`;
      element.hidden = false;
    },

    wyczysc() {
      ostatnia = '';
      transkrypcja.textContent = '';
      odnosnik.textContent = '';
      element.hidden = true;
    },
  };
}

/**
 * Zaznaczony fragment odpowiedzi albo jej całość.
 *
 * Zaznaczenie czytamy przez `getSelection`, ale nie zakładamy, że jest —
 * środowisko sprawdzianu i przeglądarka bez zaznaczenia oddają wtedy całość,
 * zamiast pustego zadania.
 */
function fragment(zrodlo: HTMLElement, calosc: string): string {
  const zaznaczenie = typeof window.getSelection === 'function' ? window.getSelection() : null;
  const tekst = (zaznaczenie?.toString() ?? '').trim();
  if (tekst !== '' && zrodlo.textContent?.includes(tekst) === true) return tekst;
  return calosc;
}
