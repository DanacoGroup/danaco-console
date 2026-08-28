import { describe, expect, it } from 'vitest';

import { Command, ComponentKind, ConfigScope, type Component } from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';
import { utworzPanelKomponentu } from './panel-komponentu';
import { utworzZmianeKomponentu } from './zmiana-komponentu';

/** Komponent własny już założony obsługuje zmianę pól i przypisanie do poziomu zasięgu, w postaci, w której rdzeń go oddaje, z domyślnymi wartościami nadpisywanymi zmianami wskazanymi w wywołaniu. */
function komponent(zmiany: Partial<Component> = {}): Component {
  return {
    id: 'komponent-1',
    kind: ComponentKind.Automations,
    name: 'Nocna automatyka',
    enabled: true,
    createdAt: 1_700_000_000_000,
    updatedAt: 1_700_000_600_000,
    ...zmiany,
  };
}

/** Kanał próbny zapamiętuje żądania wysłane w sprawdzianie i oddaje odpowiedź wskazaną osobno dla każdej komendy. */
function kanalProbny(odpowiedzi: Record<string, unknown>): {
  kanal: Kanal;
  wyslane: { komenda: string; zadanie: unknown }[];
} {
  const wyslane: { komenda: string; zadanie: unknown }[] = [];
  const kanal = {
    wyslij(komenda: string, zadanie: unknown, przyWyniku?: (wynik: Wynik<unknown>) => void) {
      wyslane.push({ komenda, zadanie });
      const tresc = odpowiedzi[komenda];
      przyWyniku?.(
        tresc === undefined
          ? { udany: false, blad: { code: 'not_found', message: 'brak wykonawcy', retryable: false } }
          : { udany: true, wynik: tresc },
      );
      return `zadanie-${String(wyslane.length)}`;
    },
    naZdarzenie: () => () => undefined,
    naDowolny: () => () => undefined,
    sesja: () => ({ id: () => 'sesja-1' }) as unknown as ReturnType<Kanal['sesja']>,
    dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
  } as unknown as Kanal;
  return { kanal, wyslane };
}

/** Panel rozwinięty wraz z odczytanym wykazem, gotowy do sprawdzania kontrolek i wysłanych żądań próby. */
async function panel(odpowiedzi: Record<string, unknown>) {
  const { kanal, wyslane } = kanalProbny(odpowiedzi);
  let odswiezono = 0;
  const zbudowany = utworzPanelKomponentu({
    kanal,
    naZmiane: () => {
      odswiezono += 1;
    },
    zrodlo: utworzZmianeKomponentu(kanal),
  });
  zbudowany.przelacz();
  await przemiel();
  return { panel: zbudowany, element: zbudowany.element, wyslane, odswiezono: () => odswiezono };
}

async function przemiel(): Promise<void> {
  for (let krok = 0; krok < 6; krok += 1) await Promise.resolve();
}

/** Kontrolka panelu o wskazanym znaczniku steru, znaleziona przez selektor atrybutu danych na elemencie. */
function ster<T extends HTMLElement>(element: HTMLElement, nazwa: string): T | null {
  return element.querySelector<T>(`[data-ster="${nazwa}"]`);
}

/** Przycisk panelu o wskazanym napisie, znaleziony przeszukaniem wszystkich przycisków danego elementu. */
function przycisk(element: HTMLElement, napis: string): HTMLButtonElement | undefined {
  return [...element.querySelectorAll('button')].find((pozycja) => pozycja.textContent === napis);
}

describe('zmiana komponentu własnego (component.update)', () => {
  it('czyta wykaz wraz z komponentami wyłączonymi — te też się zmienia', async () => {
    const { wyslane } = await panel({
      [Command.ComponentList]: { components: [komponent()] },
    });

    const wykaz = wyslane.find((pozycja) => pozycja.komenda === Command.ComponentList);
    expect(wykaz?.zadanie).toEqual({ includeDisabled: true });
  });

  it('wysyła wyłącznie pola dotknięte', async () => {
    const { element, wyslane } = await panel({
      [Command.ComponentList]: { components: [komponent()] },
      [Command.ComponentUpdate]: { component: komponent({ name: 'Automatyka poranna' }) },
    });

    const nazwa = element.querySelector<HTMLInputElement>('[aria-label="Nowa nazwa komponentu"]');
    if (nazwa !== null) nazwa.value = 'Automatyka poranna';
    przycisk(element, 'Zapisz zmianę')?.click();
    await przemiel();

    const zmiana = wyslane.find((pozycja) => pozycja.komenda === Command.ComponentUpdate);
    // Ani opisu, ani stanu czynności: pole nietknięte zostaje w rdzeniu bez zmian.
    expect(zmiana?.zadanie).toEqual({ componentId: 'komponent-1', name: 'Automatyka poranna' });
  });

  it('żądania bez ani jednej zmiany nie wysyła i mówi dlaczego', async () => {
    const { element, wyslane } = await panel({
      [Command.ComponentList]: { components: [komponent()] },
    });

    przycisk(element, 'Zapisz zmianę')?.click();
    await przemiel();

    expect(wyslane.some((pozycja) => pozycja.komenda === Command.ComponentUpdate)).toBe(false);
    expect(element.textContent ?? '').toContain('Pole puste zostawia wartość bez zmian');
  });

  it('stan czynności jedzie tylko wtedy, gdy różni się od rdzenia', async () => {
    const { element, wyslane } = await panel({
      [Command.ComponentList]: { components: [komponent({ enabled: true })] },
      [Command.ComponentUpdate]: { component: komponent({ enabled: false }) },
    });

    const czynny = element.querySelector<HTMLInputElement>('#dn-komponent-czynny');
    if (czynny !== null) czynny.checked = false;
    przycisk(element, 'Zapisz zmianę')?.click();
    await przemiel();

    const zmiana = wyslane.find((pozycja) => pozycja.komenda === Command.ComponentUpdate);
    expect(zmiana?.zadanie).toEqual({ componentId: 'komponent-1', enabled: false });
  });
});

describe('przypisanie komponentu (component.assign)', () => {
  it('mówi, z czym wiąże, zanim zwiąże', async () => {
    const { element } = await panel({
      [Command.ComponentList]: { components: [komponent()] },
      [Command.EnvironmentList]: {
        environments: [
          {
            id: 'srodowisko-1',
            code: 'praca',
            name: 'Praca',
            order: 1,
            navigationKind: 'modules',
          },
        ],
      },
    });

    const poziom = ster<HTMLSelectElement>(element, 'poziom');
    if (poziom !== null) {
      poziom.value = ConfigScope.Environment;
      poziom.dispatchEvent(new Event('change'));
    }
    await przemiel();

    const zdanie = element.textContent ?? '';
    // Nazwa komponentu i nazwa bytu, nie identyfikatory, i przed czynnością wiązania.
    expect(zdanie).toContain('Wiążesz komponent „Nocna automatyka" z: Praca (praca)');
    expect(przycisk(element, 'Przypisz komponent')?.disabled).toBe(false);
  });

  it('poziom globalny wiąże bez bytu i mówi to wprost', async () => {
    const { element, wyslane } = await panel({
      [Command.ComponentList]: { components: [komponent()] },
      [Command.ComponentAssign]: { component: komponent(), assigned: true },
    });

    expect(element.textContent ?? '').toContain('z poziomem globalnym — bez bytu');
    przycisk(element, 'Przypisz komponent')?.click();
    await przemiel();

    const przypisanie = wyslane.find((pozycja) => pozycja.komenda === Command.ComponentAssign);
    expect(przypisanie?.zadanie).toEqual({
      componentId: 'komponent-1',
      scope: ConfigScope.Global,
    });
  });

  it('powtórzone przypisanie nie udaje czynności, której rdzeń nie wykonał', async () => {
    const { element } = await panel({
      [Command.ComponentList]: { components: [komponent()] },
      [Command.ComponentAssign]: { component: komponent(), assigned: false },
    });

    przycisk(element, 'Przypisz komponent')?.click();
    await przemiel();

    const zdanie = element.textContent ?? '';
    expect(zdanie).toContain('Nic się nie zmieniło');
    expect(zdanie).toContain('przypisanie nie doszło do skutku');
  });

  it('mówi, po czyjej stronie brak poziomów, których rdzeń nie przyjmuje', async () => {
    const { element } = await panel({
      [Command.ComponentList]: { components: [komponent()] },
    });

    expect(element.textContent ?? '').toContain('Brak jest po stronie RDZENIA');
  });

  it('bez ani jednego komponentu nie wiąże i mówi, czego brakuje', async () => {
    const { element } = await panel({
      [Command.ComponentList]: { components: [] },
    });

    expect(element.textContent ?? '').toContain('Nie ma czego wiązać');
    expect(przycisk(element, 'Przypisz komponent')?.disabled).toBe(true);
  });
});
