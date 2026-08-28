import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { nazwaZasobu } from './karta-zasobu';
import type { StanDesignu } from './stan-designu';

/**
 * Udostępnia dwie czynności na zasobie wskazanym w wykazie: oznaczenie ulubionym
 * i usunięcie bez dodatkowego potwierdzenia, z paskiem odpowiedzi po wykonaniu.
 */
export interface CzynnosciZasobu {
  element: HTMLElement;
  /** Przepisuje napis przełącznika pod zasób wskazany; woła je odświeżenie panelu. */
  odswiez(): void;
}

export function utworzCzynnosciZasobu(stan: StanDesignu): CzynnosciZasobu {
  const ulubiony = przycisk('Oznacz ulubionym', 'dn-btn dn-btn--sm dn-btn--zarys');
  const usun = przycisk('Usuń zasób', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const pasek = document.createElement('div');
  pasek.className = 'md-czynnosci__pasek';
  pasek.append(ulubiony, usun);

  const element = document.createElement('div');
  element.className = 'md-czynnosci';
  element.append(pasek, odpowiedz.element);

  ulubiony.addEventListener('click', () => void przestawUlubiony());
  usun.addEventListener('click', () => void usunWskazany());

  /** Przestawia oznaczenie ulubionego zasobu; stan liczy się z zasobu, nie z napisu przycisku. */
  async function przestawUlubiony(): Promise<void> {
    const zasob = stan.wybrany();
    if (zasob === null) {
      odpowiedz.pokaz('Wskaż zasób w wykazie — oznaczenie dotyczy jednego zasobu.', false);
      return;
    }
    const docelowy = zasob.favorite !== true;
    odpowiedz.pokaz(
      `${docelowy ? 'Oznaczanie' : 'Zdejmowanie oznaczenia'} „${nazwaZasobu(zasob)}"…`,
      true,
    );
    const wynik = await stan.czuwanie.prowadz(
      'oznaczenie ulubionego',
      stan.zrodlo.ustawUlubiony(zasob.id, docelowy),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Oznaczenie ulubionego', wynik.blad), false);
      return;
    }
    // Rdzeń oddaje zasób po zapisie: porównanie zamówienia z odpowiedzią jest tanie i wiarygodne.
    const oddany = wynik.wynik.asset;
    stan.wchlon(oddany);
    if ((oddany.favorite === true) !== docelowy) {
      odpowiedz.pokaz(
        `Zamówiono ulubione = ${docelowy}, a rdzeń oddał zasób z ulubione = ` +
          `${oddany.favorite === true}. Zmiany nie potwierdzam.`,
        false,
      );
      return;
    }
    odpowiedz.pokaz(
      docelowy
        ? `Zasób „${nazwaZasobu(oddany)}" jest ulubiony — przełącznik „Tylko ulubione" w filtrze ` +
          'zawęża teraz wykaz także do niego.'
        : `Zasób „${nazwaZasobu(oddany)}" przestał być ulubiony.`,
      true,
    );
  }

  /** Usuwa zasób wskazany, bez bramki potwierdzającej przed wykonaniem czynności. */
  async function usunWskazany(): Promise<void> {
    const zasob = stan.wybrany();
    if (zasob === null) {
      odpowiedz.pokaz('Wskaż zasób w wykazie — usunięcie dotyczy jednego zasobu.', false);
      return;
    }
    const nazwa = nazwaZasobu(zasob);
    odpowiedz.pokaz(`Usuwanie zasobu „${nazwa}"…`, true);
    const wynik = await stan.czuwanie.prowadz(
      'usunięcie zasobu',
      stan.zrodlo.usunZasob(zasob.id),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        // Cisza kanału nie jest odmową: rdzeń mógł zasób skasować, a odpowiedź zginąć w sieci.
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Usunięcie zasobu', wynik.blad), false);
      return;
    }
    if (!wynik.wynik.removed) {
      // Zasób był w wykazie, a rdzeń go nie zna — wykaz był nieaktualny; usuwamy wpis z listy.
      stan.zdejmij(zasob.id);
      odpowiedz.pokaz(
        `Rdzeń nie zna zasobu „${nazwa}" (${zasob.id}) — niczego nie usunął. Wykaz był ` +
          'nieaktualny i pozycja została z niego zdjęta.',
        false,
      );
      return;
    }
    stan.zdejmij(zasob.id);
    odpowiedz.pokaz(`Rdzeń usunął zasób „${nazwa}" (${zasob.id}).`, true);
  }

  return {
    element,

    odswiez() {
      const zasob = stan.wybrany();
      // Napis przycisku opisuje czynność do wykonania, a nie bieżący stan zasobu.
      ulubiony.textContent =
        zasob?.favorite === true ? 'Zdejmij oznaczenie ulubionego' : 'Oznacz ulubionym';
    },
  };
}
