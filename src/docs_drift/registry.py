import re
from pathlib import Path
from typing import Any

import yaml  # type: ignore[import-untyped]  # types-PyYAML is not a declared dependency; see README
from pydantic import BaseModel, ConfigDict, ValidationError, field_validator, model_validator
from pydantic_core import ErrorDetails

REPO_PATTERN = re.compile(r"^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$")
CREDENTIAL_WORDS = {"token", "password", "secret", "key"}


def _check_relative_path(value: str) -> str:
    if value.startswith("/") or ".." in value.split("/"):
        raise ValueError(f"path '{value}' must stay within the repository: use a relative path without '..'")
    return value


class CodeSource(BaseModel):
    model_config = ConfigDict(extra="forbid")

    repo: str
    ref: str

    @model_validator(mode="before")
    @classmethod
    def _reject_credential_fields(cls, data: Any) -> Any:  # raw YAML is untyped until validated
        return _reject_credential_keys(data)

    @field_validator("repo")
    @classmethod
    def _check_repo_format(cls, value: str) -> str:
        if not REPO_PATTERN.match(value):
            raise ValueError(f"repo '{value}' must use owner/name form")
        return value

    @field_validator("ref")
    @classmethod
    def _check_ref_not_empty(cls, value: str) -> str:
        if not value.strip():
            raise ValueError("ref must not be empty")
        return value


class Project(BaseModel):
    model_config = ConfigDict(extra="forbid")

    docs: str
    code: list[CodeSource] = []
    openapi: str | None = None

    @model_validator(mode="before")
    @classmethod
    def _reject_credential_fields(cls, data: Any) -> Any:  # raw YAML is untyped until validated
        return _reject_credential_keys(data)

    @field_validator("docs")
    @classmethod
    def _check_docs_path(cls, value: str) -> str:
        return _check_relative_path(value)

    @field_validator("openapi")
    @classmethod
    def _check_openapi_path(cls, value: str | None) -> str | None:
        return None if value is None else _check_relative_path(value)

    @field_validator("code")
    @classmethod
    def _check_code_not_empty(cls, value: list[CodeSource]) -> list[CodeSource]:
        if not value:
            raise ValueError("at least one code source is required")
        return value


class Registry(BaseModel):
    model_config = ConfigDict(extra="forbid")

    projects: dict[str, Project]


class RegistryError(Exception):
    def __init__(self, errors: list[str]) -> None:
        super().__init__("\n".join(errors))
        self.errors = errors


def load_registry(path: Path) -> Registry:
    try:
        raw = yaml.safe_load(path.read_text(encoding="utf-8"))
    except yaml.YAMLError as exc:
        raise RegistryError([f"{path}: invalid YAML: {exc}"]) from exc
    try:
        return Registry.model_validate(raw)
    except ValidationError as exc:
        raise RegistryError([_format_error(error) for error in exc.errors()]) from exc


def _reject_credential_keys(data: Any) -> Any:
    if isinstance(data, dict):
        for key in data:
            if isinstance(key, str) and CREDENTIAL_WORDS.intersection(re.split(r"[_-]", key.lower())):
                raise ValueError(f"field '{key}' is not allowed; credentials belong in the secret store")
    return data


def _format_error(error: ErrorDetails) -> str:
    loc = list(error["loc"])
    message = str(error["msg"]).removeprefix("Value error, ")
    prefix = ""
    if len(loc) >= 2 and loc[0] == "projects":
        prefix = f"{loc[1]}: "
        loc = loc[2:]
    field = _format_path(loc)
    return f"{prefix}{field}: {message}" if field else f"{prefix}{message}"


def _format_path(loc: list[Any]) -> str:
    path = ""
    for part in loc:
        path += f"[{part}]" if isinstance(part, int) else f".{part}"
    return path.removeprefix(".")
