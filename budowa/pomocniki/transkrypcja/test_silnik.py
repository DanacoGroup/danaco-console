#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Sprawdziany silnik.py — atrapa faster_whisper zamiast prawdziwej biblioteki."""

from __future__ import annotations

import sys
import types
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))

import silnik  # noqa: E402


class _AtrapaModel:
    """Przechwytuje kwargs wołania WhisperModel zamiast ładować wagi."""

    ostatnie_kwargs: dict = {}

    def __init__(self, model, **kwargs):
        _AtrapaModel.ostatnie_kwargs = kwargs

    def transcribe(self, sciezka, language=None):
        info = types.SimpleNamespace(duration=1.5)
        return iter([]), info


def test_transkrybuj_woła_whispermodel_z_local_files_only(monkeypatch):
    atrapa = types.ModuleType("faster_whisper")
    atrapa.WhisperModel = _AtrapaModel
    monkeypatch.setitem(sys.modules, "faster_whisper", atrapa)

    silnik.transkrybuj(Path("nagranie.wav"), "small", "pl")

    assert _AtrapaModel.ostatnie_kwargs.get("local_files_only") is True
