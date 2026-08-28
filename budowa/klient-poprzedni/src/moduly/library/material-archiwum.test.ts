import { describe, expect, it } from 'vitest';

import {
  Command,
  DesignAssetKind,
  MediaOperationKind,
  type DesignAsset,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { utworzPanelPakowania } from './archiwum-pakowanie';
import { utworzNarzedziaMaterialu } from './material-narzedzia';
import { utworzPanelMaterialu } from './material-panel';

// Materiał i archiwum: trzy czynności arsenału na drodze z okna, ścieżce z dysku i nazwanym wyniku.

/** Zasób w postaci, w której rdzeń oddaje wynik czynności arsenału, gotowy do sprawdzenia jego pól w sprawdzianie. */
function zasob(zmiany: Partial<DesignAsset> = {}): DesignAsset {
  return {
    id: 'zasob-1',
    windowId: '',
    kind: DesignAssetKind.Video,
    name: 'nagranie-po-cieciu',
    format: 'mp4',
    createdAt: 1_700_000_000_000,
    ...zmiany,
  };
}

/** Kanał próbny zapamiętuje wysłane żądania i oddaje odpowiedź wskazaną dla danej komendy, imitując rdzeń. */
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

/** Wpisuje wartość w pole panelu o wskazanej etykiecie dostępności, symulując wpis wykonany przez człowieka. */
function wpisz(element: HTMLElement, znacznik: string, wartosc: string): void {
  const pole = element.querySelector<HTMLInputElement>(`[data-pole="${znacznik}"]`);
  expect(pole, `pole „${znacznik}" musi stać w panelu`).not.toBeNull();
  if (pole !== null) pole.value = wartosc;
}

/** Naciska przycisk czynności i oddaje treść wiersza odpowiedzi widocznego po jej wykonaniu w tym oknie. */
async function nacisnij(element: HTMLElement, czynnosc: string): Promise<string> {
  const przycisk = element.querySelector<HTMLButtonElement>(`[data-czynnosc="${czynnosc}"]`);
  expect(przycisk, `czynność „${czynnosc}" musi mieć przycisk`).not.toBeNull();
  przycisk?.click();
  await Promise.resolve();
  await Promise.resolve();
  return element.querySelector('.dm-odpowiedz')?.textContent ?? '';
}

describe('rozpoznanie materiału (media.inspect)', () => {
  it('woła rdzeń ścieżką z dysku i nazywa kontener, czas oraz strumienie', async () => {
    const { kanal, wyslane } = kanalProbny({
      [Command.MediaInspect]: {
        durationMs: 61_000,
        streams: 'video/h264, audio/aac',
        format: 'mov,mp4,m4a',
        sizeBytes: 8_400_000,
      },
    });
    const panel = utworzPanelMaterialu(utworzNarzedziaMaterialu(kanal));

    wpisz(panel.element, 'sciezka', '/dom/operator/nagranie.mp4');
    const zdanie = await nacisnij(panel.element, 'material-rozpoznanie');

    expect(wyslane).toHaveLength(1);
    expect(wyslane[0]?.komenda).toBe(Command.MediaInspect);
    // Droga zasobu biblioteki wróciłaby odmową, więc żądanie nie może nieść identyfikatora zasobu.
    expect(wyslane[0]?.zadanie).toEqual({ sourcePath: '/dom/operator/nagranie.mp4' });
    expect(zdanie).toContain('mov,mp4,m4a');
    expect(zdanie).toContain('61000 ms');
    expect(zdanie).toContain('video/h264, audio/aac');
  });

  it('nazywa brak zapisanego czasu trwania zamiast pokazać zero milisekund', async () => {
    const { kanal } = kanalProbny({
      [Command.MediaInspect]: { durationMs: 0, streams: 'audio/opus', format: 'ogg', sizeBytes: 12 },
    });
    const panel = utworzPanelMaterialu(utworzNarzedziaMaterialu(kanal));

    wpisz(panel.element, 'sciezka', '/dom/operator/strumien.ogg');
    const zdanie = await nacisnij(panel.element, 'material-rozpoznanie');

    expect(zdanie).toContain('czasu trwania materiał nie ma zapisanego');
    expect(zdanie).not.toContain('czas trwania 0 ms');
  });

  it('bez ścieżki nie woła rdzenia i mówi, czego brakuje', async () => {
    const { kanal, wyslane } = kanalProbny({});
    const panel = utworzPanelMaterialu(utworzNarzedziaMaterialu(kanal));

    const zdanie = await nacisnij(panel.element, 'material-rozpoznanie');

    expect(wyslane).toHaveLength(0);
    expect(zdanie).toContain('nie ma czego zmierzyć');
  });
});

describe('przetworzenie materiału (media.transcode)', () => {
  it('oddaje wynik nazwany i mówi, gdzie leżą bajty', async () => {
    const { kanal, wyslane } = kanalProbny({
      [Command.MediaTranscode]: { asset: zasob(), durationMs: 4_000 },
    });
    const panel = utworzPanelMaterialu(utworzNarzedziaMaterialu(kanal));

    wpisz(panel.element, 'sciezka', '/dom/operator/nagranie.mov');
    wpisz(panel.element, 'od-ms', '1000');
    wpisz(panel.element, 'do-ms', '5000');
    const wybor = panel.element.querySelector<HTMLSelectElement>('[data-pole="operacja"]');
    if (wybor !== null) wybor.value = MediaOperationKind.Trim;

    const zdanie = await nacisnij(panel.element, 'material-przetworzenie');

    expect(wyslane[0]?.zadanie).toEqual({
      sourcePath: '/dom/operator/nagranie.mov',
      operation: MediaOperationKind.Trim,
      startMs: 1000,
      endMs: 5000,
    });
    // Wynik NAZWANY: nazwa własna, format i identyfikator zasobu.
    expect(zdanie).toContain('nagranie-po-cieciu');
    expect(zdanie).toContain('zasob-1');
    expect(zdanie).toContain('magazynie zasobów');
  });

  it('pole puste nie jedzie jako zero', async () => {
    const { kanal, wyslane } = kanalProbny({
      [Command.MediaTranscode]: { asset: zasob() },
    });
    const panel = utworzPanelMaterialu(utworzNarzedziaMaterialu(kanal));

    wpisz(panel.element, 'sciezka', '/dom/operator/nagranie.mov');
    await nacisnij(panel.element, 'material-przetworzenie');

    const zadanie = wyslane[0]?.zadanie as Record<string, unknown>;
    expect(zadanie).not.toHaveProperty('startMs');
    expect(zadanie).not.toHaveProperty('width');
  });

  it('odmowę rdzenia pokazuje wprost, a nie ciszą', async () => {
    const { kanal } = kanalProbny({});
    const panel = utworzPanelMaterialu(utworzNarzedziaMaterialu(kanal));

    wpisz(panel.element, 'sciezka', '/dom/operator/nagranie.mov');
    const zdanie = await nacisnij(panel.element, 'material-przetworzenie');

    expect(zdanie).toContain('Przetworzenie materiału');
    expect(zdanie.length).toBeGreaterThan(20);
  });
});

describe('spakowanie archiwum (archive.pack)', () => {
  it('pakuje ścieżkę i nazywa wynik wraz z liczbą pozycji', async () => {
    const { kanal, wyslane } = kanalProbny({
      [Command.ArchivePack]: {
        asset: zasob({ id: 'zasob-7', name: 'wydanie', format: 'zip' }),
        entries: 12,
        sizeBytes: 4096,
      },
    });
    const panel = utworzPanelPakowania(utworzNarzedziaMaterialu(kanal));

    wpisz(panel.element, 'sciezka', '/dom/operator/praca');
    const zdanie = await nacisnij(panel.element, 'archiwum-pakowanie');

    expect(wyslane[0]?.komenda).toBe(Command.ArchivePack);
    expect(wyslane[0]?.zadanie).toEqual({ sourcePath: '/dom/operator/praca' });
    expect(zdanie).toContain('wydanie');
    expect(zdanie).toContain('zasob-7');
    expect(zdanie).toContain('Pozycji w archiwum: 12');
  });

  it('archiwum bez pozycji nazywa pustkę, zamiast milczeć o niej', async () => {
    const { kanal } = kanalProbny({
      [Command.ArchivePack]: { asset: zasob({ format: 'zip' }), entries: 0, sizeBytes: 22 },
    });
    const panel = utworzPanelPakowania(utworzNarzedziaMaterialu(kanal));

    wpisz(panel.element, 'sciezka', '/dom/operator/pusty');
    const zdanie = await nacisnij(panel.element, 'archiwum-pakowanie');

    expect(zdanie).toContain('ani jednej pozycji');
  });

  it('postać i nazwę wysyła tylko wtedy, gdy Operator je wskazał', async () => {
    const { kanal, wyslane } = kanalProbny({
      [Command.ArchivePack]: { asset: zasob(), entries: 1, sizeBytes: 10 },
    });
    const panel = utworzPanelPakowania(utworzNarzedziaMaterialu(kanal));

    wpisz(panel.element, 'sciezka', '/dom/operator/praca');
    const postac = panel.element.querySelector<HTMLSelectElement>('[data-pole="postac"]');
    if (postac !== null) postac.value = '7z';
    wpisz(panel.element, 'nazwa', 'wydanie-sierpniowe');
    await nacisnij(panel.element, 'archiwum-pakowanie');

    expect(wyslane[0]?.zadanie).toEqual({
      sourcePath: '/dom/operator/praca',
      format: '7z',
      name: 'wydanie-sierpniowe',
    });
  });
});
