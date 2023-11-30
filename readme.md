<p align="center">
  <picture>
    <source srcset="imgs/mship_gopher.png" media="(prefers-color-scheme: dark)" height="150">
    <img src="imgs/mship_gopher_dark.png" alt="Mship" height="150">
  </picture>
</p>
<p align="center">Tool to archive RPM packages and attest to their authenticity</p>
<hr />

**This code is being moved from the Rocky Linux build system monorepo and thus documentation and
instructions are not fully updated yet.**

### Development

Using the taskrunner2 target is sufficient. `bazel run //tools/mothership`.

This will watch for changes in mship_admin_server, mship_server and both UIs.

The target will also start Dex IDP and Temporal if not already running.

### Fun fact

Mship (or Mothership) was created at the RESF for Rocky Linux after the
announcement regarding RHEL public source code access.
The name is a play on the fact that the RHEL source code is being imported
by the mother ship and that the downstream builders are the smaller ships
being served by it.
