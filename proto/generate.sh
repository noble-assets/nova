cd proto
buf generate --template buf.gen.gogo.yaml
buf generate --template buf.gen.pulsar.yaml
cd ..

cp -r github.com/noble-assets/nova/* ./
cp -r api/nova/* api/
find api/ -type f -name "*.go" -exec sed -i 's|github.com/noble-assets/nova/api/nova|github.com/noble-assets/nova/api|g' {} +

rm -rf github.com
rm -rf api/nova
rm -rf nova
