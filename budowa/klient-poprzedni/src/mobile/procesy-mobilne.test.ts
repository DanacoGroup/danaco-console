import { describe, expect, it } from 'vitest';

import {
  Command,
  MobileProcessControl,
  ProgressStatus,
  type MobileProcess,
} from '../../../shared/contract';
import { utworzEkranProcesow } from './ekran-procesow';
import type { Kanal, Wynik } from '../protokol/kanal';
import { utworzZrodloProcesowMobilnych } from './procesy-mobilne';

// Sprawdzian pilnuje drogi z okna, filtra po stronie rdzenia i nieodwracalności czynności.

/** Proces mobilny w postaci, w jakiej dokładnie rdzeń go oddaje odpowiedzią komendy `mobile.process.list`. */
function proces(zmiany: Partial<MobileProcess> = {}): MobileProcess {
  return {
    id: 'proces-1',
    label: 'Budowa wydania',
    status: ProgressStatus.Running,
    completion: 40,
    ...zmiany,
  };
}

/** Kanał próbny do sprawdzianu: zapamiętuje wysłane żądania i oddaje odpowiedź wskazaną per komenda testu. */
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
    sesja: () => ({ id: () => 'urzadzenie-1' }) as unknown as ReturnType<Kanal['sesja']>,
    dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
  } as unknown as Kanal;
  return { kanal, wyslane };
}

function ekran(odpowiedzi: Record<string, unknown>) {
  const { kanal, wyslane } = kanalProbny(odpowiedzi);
  const zbudowany = utworzEkranProcesow({
    kanal,
    zrodlo: utworzZrodloProcesowMobilnych(kanal),
  });
  return { ekran: zbudowany, element: zbudowany.element, wyslane };
}

/** Przycisk czynności sterowania przy wierszu procesu, znaleziony po identyfikatorze i rodzaju czynności. */
function przyciskCzynnosci(
  element: HTMLElement,
  idProcesu: string,
  czynnosc: MobileProcessControl,
): HTMLButtonElement | null {
  return element.querySelector<HTMLButtonElement>(
    `[data-proces="${idProcesu}"] [data-czynnosc="${czynnosc}"]`,
  );
}

describe('wykaz procesów (mobile.process.list)', () => {
  it('woła rdzeń z urządzeniem wziętym z karty sesji kanału', async () => {
    const { ekran: zbudowany, wyslane } = ekran({
      [Command.MobileProcessList]: { processes: [proces()] },
    });

    await zbudowany.odswiez();

    expect(wyslane).toHaveLength(1);
    expect(wyslane[0]?.komenda).toBe(Command.MobileProcessList);
    expect(wyslane[0]?.zadanie).toEqual({ deviceId: 'urzadzenie-1' });
    expect(zbudowany.procesy()).toHaveLength(1);
  });

  it('pokazuje proces wraz ze stanem i stopniem ukończenia', async () => {
    const { ekran: zbudowany, element } = ekran({
      [Command.MobileProcessList]: { processes: [proces()] },
    });

    await zbudowany.odswiez();

    const tresc = element.textContent ?? '';
    expect(tresc).toContain('Budowa wydania');
    expect(tresc).toContain('W biegu');
    expect(tresc).toContain('ukończone 40%');
  });

  it('brak stopnia ukończenia nazywa, zamiast pokazać zero procent', async () => {
    const { ekran: zbudowany, element } = ekran({
      [Command.MobileProcessList]: {
        processes: [{ id: 'proces-2', label: 'Oczekujące', status: ProgressStatus.Pending }],
      },
    });

    await zbudowany.odswiez();

    expect(element.textContent ?? '').toContain('stopnia ukończenia rdzeń nie podaje');
    expect(element.textContent ?? '').not.toContain('ukończone 0%');
  });

  it('filtr stanu zawęża po stronie rdzenia, a nie w oknie', async () => {
    const { ekran: zbudowany, element, wyslane } = ekran({
      [Command.MobileProcessList]: { processes: [proces({ status: ProgressStatus.Paused })] },
    });
    await zbudowany.odswiez();

    const filtr = element.querySelector<HTMLSelectElement>('select');
    expect(filtr).not.toBeNull();
    if (filtr !== null) {
      filtr.value = ProgressStatus.Paused;
      filtr.dispatchEvent(new Event('change'));
    }
    await Promise.resolve();
    await Promise.resolve();

    const ostatnie = wyslane[wyslane.length - 1];
    expect(ostatnie?.zadanie).toEqual({
      deviceId: 'urzadzenie-1',
      status: ProgressStatus.Paused,
    });
  });

  it('pusty wykaz ma zdanie, nie puste miejsce', async () => {
    const { ekran: zbudowany, element } = ekran({
      [Command.MobileProcessList]: { processes: [] },
    });

    await zbudowany.odswiez();

    expect(element.textContent ?? '').toContain('ani jednego procesu');
  });

  it('odmowę rdzenia pokazuje wprost i nie zostawia wykazu zastanego', async () => {
    const { ekran: zbudowany, element } = ekran({});

    await zbudowany.odswiez();

    expect(element.textContent ?? '').toContain('Wykaz procesów platformy');
    expect(zbudowany.procesy()).toHaveLength(0);
  });
});

describe('sterowanie procesem (mobile.process.control)', () => {
  it('czynność odwracalna idzie do rdzenia od razu', async () => {
    const { ekran: zbudowany, element, wyslane } = ekran({
      [Command.MobileProcessList]: { processes: [proces()] },
      [Command.MobileProcessControl]: { process: proces({ status: ProgressStatus.Paused }) },
    });
    await zbudowany.odswiez();

    przyciskCzynnosci(element, 'proces-1', MobileProcessControl.Pause)?.click();
    await Promise.resolve();
    await Promise.resolve();

    const sterowanie = wyslane.filter((pozycja) => pozycja.komenda === Command.MobileProcessControl);
    expect(sterowanie).toHaveLength(1);
    expect(sterowanie[0]?.zadanie).toEqual({
      processId: 'proces-1',
      control: MobileProcessControl.Pause,
      deviceId: 'urzadzenie-1',
    });
  });

  it('zatrzymanie mówi o nieodwracalności PRZED wykonaniem i nie woła rdzenia', async () => {
    const { ekran: zbudowany, element, wyslane } = ekran({
      [Command.MobileProcessList]: { processes: [proces()] },
      [Command.MobileProcessControl]: { process: proces({ status: ProgressStatus.Stopped }) },
    });
    await zbudowany.odswiez();

    const przycisk = przyciskCzynnosci(element, 'proces-1', MobileProcessControl.Stop);
    przycisk?.click();
    await Promise.resolve();

    // Pierwsze dotknięcie uzbraja i ostrzega — do rdzenia nie idzie nic.
    expect(wyslane.filter((p) => p.komenda === Command.MobileProcessControl)).toHaveLength(0);
    const tresc = element.textContent ?? '';
    expect(tresc).toContain('Czynność nieodwracalna');
    expect(tresc).toContain('nie da się go wznowić');
    expect(przycisk?.dataset['uzbrojona']).toBe('tak');
    expect(przycisk?.textContent).toContain('dotknij ponownie');
  });

  it('drugie dotknięcie zatrzymania wykonuje je i odczytuje wykaz na nowo', async () => {
    const { ekran: zbudowany, element, wyslane } = ekran({
      [Command.MobileProcessList]: { processes: [proces()] },
      [Command.MobileProcessControl]: { process: proces({ status: ProgressStatus.Stopped }) },
    });
    await zbudowany.odswiez();

    const przycisk = przyciskCzynnosci(element, 'proces-1', MobileProcessControl.Stop);
    przycisk?.click();
    await Promise.resolve();
    przyciskCzynnosci(element, 'proces-1', MobileProcessControl.Stop)?.click();
    await Promise.resolve();
    await Promise.resolve();
    await Promise.resolve();

    const sterowanie = wyslane.filter((pozycja) => pozycja.komenda === Command.MobileProcessControl);
    expect(sterowanie).toHaveLength(1);
    expect(sterowanie[0]?.zadanie).toEqual({
      processId: 'proces-1',
      control: MobileProcessControl.Stop,
      deviceId: 'urzadzenie-1',
    });
    // Rodzina nie ma zdarzenia własnego, więc wykaz czyta się po sterowaniu.
    expect(wyslane.filter((p) => p.komenda === Command.MobileProcessList).length).toBeGreaterThan(1);
  });

  it('uzbrojenie gaśnie po odczycie wykazu — przycisk nie zostaje gotowy do strzału', async () => {
    const { ekran: zbudowany, element, wyslane } = ekran({
      [Command.MobileProcessList]: { processes: [proces()] },
      [Command.MobileProcessControl]: { process: proces() },
    });
    await zbudowany.odswiez();

    przyciskCzynnosci(element, 'proces-1', MobileProcessControl.Stop)?.click();
    await Promise.resolve();
    await zbudowany.odswiez();
    przyciskCzynnosci(element, 'proces-1', MobileProcessControl.Stop)?.click();
    await Promise.resolve();

    expect(wyslane.filter((p) => p.komenda === Command.MobileProcessControl)).toHaveLength(0);
  });
});
