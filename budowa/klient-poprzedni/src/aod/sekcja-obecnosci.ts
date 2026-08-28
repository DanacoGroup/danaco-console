import type { AodStatus } from '../../../shared/contract';
import { opisOdmowyAod } from './odmowy-aod';
import { utworzAkapit, utworzPodtytul } from './pola-wykazu';
import type { ZrodloAod } from './zrodlo-komend';

/** Rejestr obecności procesów przypiętych do nakładki, rysujący odpowiedź rdzenia, a nie przewidywanie klienta. */
export interface SekcjaObecnosci {
  element: HTMLElement;
  /** Nanosi wykaz przypięć i licznik procesów w biegu ze stanu nakładki. */
  pokaz(status: AodStatus): void;
  // Nanosi odmowę odczytu stanu; rejestr pusty i rejestr nieznany to dwa różne stany.
  odmowa(zdanie: string): void;
}

export interface OpisSekcjiObecnosci {
  zrodlo: ZrodloAod;
  /** Krótkie potwierdzenie czynności na pasku okna. */
  zamelduj(zdanie: string, udane: boolean): void;
}

export function utworzSekcjeObecnosci(opis: OpisSekcjiObecnosci): SekcjaObecnosci {
  const element = document.createElement('section');
  element.className = 'ao-sekcja';

  const licznik = document.createElement('p');
  licznik.className = 'ao-obecnosc__licznik';

  const miejsce = document.createElement('div');
  miejsce.className = 'ao-sekcja__tresc';

  const pole = document.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole-kontrolka ao-obecnosc__pole';
  pole.id = 'ao-obecnosc-proces';
  pole.placeholder = 'identyfikator procesu';

  const etykieta = document.createElement('label');
  etykieta.className = 'dn-pole-etykieta';
  etykieta.htmlFor = pole.id;
  etykieta.textContent = 'Przypnij proces do obserwacji';

  const przypnij = document.createElement('button');
  przypnij.type = 'submit';
  przypnij.className = 'dn-btn dn-btn--zarys';
  przypnij.textContent = 'Przypnij';

  const formularz = document.createElement('form');
  formularz.className = 'ao-formularz';
  formularz.append(etykieta, pole, przypnij);
  formularz.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    void przypnijProces();
  });

  element.append(
    utworzPodtytul('Rejestr obecności — procesy przypięte (aod.observe.attach / detach)'),
    licznik,
    miejsce,
    formularz,
  );

  /** Ostatni wykaz naniesiony na sekcję — podstawa rysowania po czynności. */
  let przypiete: readonly string[] = [];

  function narysuj(): void {
    if (przypiete.length === 0) {
      miejsce.replaceChildren(
        utworzAkapit(
          'ao-pusto',
          'Nakładka nie obserwuje żadnego procesu. Pusty rejestr jest stanem poprawnym — ' +
            'przypnij proces polem poniżej.',
        ),
      );
      return;
    }

    const lista = document.createElement('ul');
    lista.className = 'ao-obecnosc';
    for (const idProcesu of przypiete) {
      lista.append(zbudujWiersz(idProcesu));
    }
    miejsce.replaceChildren(lista);
  }

  function zbudujWiersz(idProcesu: string): HTMLLIElement {
    const wiersz = document.createElement('li');
    wiersz.className = 'ao-obecnosc__wiersz';

    const nazwa = document.createElement('span');
    nazwa.className = 'ao-obecnosc__id';
    nazwa.textContent = idProcesu;

    const odepnij = document.createElement('button');
    odepnij.type = 'button';
    odepnij.className = 'dn-btn dn-btn--zarys ao-obecnosc__przycisk';
    odepnij.textContent = 'Odepnij';
    // Bez potwierdzenia i bez wyszarzania — czynność idzie do rdzenia od razu.
    odepnij.addEventListener('click', () => void odepnijProces(idProcesu));

    wiersz.append(nazwa, odepnij);
    return wiersz;
  }

  async function przypnijProces(): Promise<void> {
    const idProcesu = pole.value.trim();
    const wynik = await opis.zrodlo.przypnij(idProcesu);

    if (!wynik.udany || wynik.wynik === undefined) {
      opis.zamelduj(opisOdmowyAod('Przypięcie procesu', 'przypiecie', wynik.blad), false);
      return;
    }

    przypiete = wynik.wynik.attachedProcessIds;
    narysuj();
    pole.value = '';
    opis.zamelduj(
      `Proces ${idProcesu} przypięty; nakładka obserwuje ${przypiete.length} proces(ów).`,
      true,
    );
  }

  async function odepnijProces(idProcesu: string): Promise<void> {
    const wynik = await opis.zrodlo.odepnij(idProcesu);

    if (!wynik.udany || wynik.wynik === undefined) {
      opis.zamelduj(opisOdmowyAod('Odpięcie procesu', 'odpiecie', wynik.blad), false);
      return;
    }

    przypiete = wynik.wynik.attachedProcessIds;
    narysuj();
    opis.zamelduj(
      `Proces ${idProcesu} odpięty; nakładka obserwuje ${przypiete.length} proces(ów).`,
      true,
    );
  }

  return {
    element,

    pokaz(status) {
      przypiete = status.attachedProcessIds ?? [];
      licznik.textContent = `Procesy w biegu widziane przez rdzeń: ${status.runningProcessCount}.`;
      narysuj();
    },

    odmowa(zdanie) {
      przypiete = [];
      licznik.textContent = 'Liczba procesów w biegu nieznana — rdzeń nie odpowiedział.';
      miejsce.replaceChildren(utworzAkapit('ao-odmowa', zdanie));
    },
  };
}
